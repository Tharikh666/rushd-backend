package prayer

import (
	"fmt"
	"time"

	"github.com/mnadev/adhango/pkg/calc"
	"github.com/mnadev/adhango/pkg/data"
	"github.com/mnadev/adhango/pkg/util"
)

// Calculator calculates daily Islamic prayer times.
type Calculator struct{}

// NewCalculator creates a new prayer calculator.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Calculate calculates prayer times for a specific location and date.
func (c *Calculator) Calculate(
	latitude float64,
	longitude float64,
	date time.Time,
	timezone string,
	settings CalculationSettings,
) (*PrayerSchedule, error) {

	// ------------------------------------------------------------
	// Validate input
	// ------------------------------------------------------------

	if err := validateCoordinates(latitude, longitude); err != nil {
		return nil, err
	}

	location, err := validateTimezone(timezone)
	if err != nil {
		return nil, err
	}

	if err := validateDate(date); err != nil {
		return nil, err
	}

	if err := validateCalculationSettings(settings); err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// Normalize date to the requested timezone.
	//
	// The date supplied by the client represents a local calendar
	// date, not a UTC date.
	// ------------------------------------------------------------

	localDate := time.Date(
		date.In(location).Year(),
		date.In(location).Month(),
		date.In(location).Day(),
		0,
		0,
		0,
		0,
		location,
	)

	// ------------------------------------------------------------
	// Coordinates
	// ------------------------------------------------------------

	coordinates, err := util.NewCoordinates(
		latitude,
		longitude,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create coordinates: %w",
			err,
		)
	}

	// ------------------------------------------------------------
	// Date components
	// ------------------------------------------------------------

	dateComponents := data.NewDateComponents(localDate)

	// ------------------------------------------------------------
	// Calculation method
	// ------------------------------------------------------------

	method, err := toAdhanCalculationMethod(
		settings.Method,
	)

	if err != nil {
		return nil, err
	}

	// Get the default parameters associated with the selected
	// calculation convention.
	parameters := calc.GetMethodParameters(method)

	if parameters == nil {
		return nil, fmt.Errorf(
			"no calculation parameters found for method %q",
			settings.Method,
		)
	}

	// ------------------------------------------------------------
	// Madhab / Asr calculation
	// ------------------------------------------------------------

	madhab, err := toAdhanMadhab(settings.Madhab)

	if err != nil {
		return nil, err
	}

	parameters.Madhab = madhab

	// ------------------------------------------------------------
	// High latitude rule
	// ------------------------------------------------------------

	highLatitudeRule, err := toAdhanHighLatitudeRule(
		settings.HighLatitudeRule,
	)

	if err != nil {
		return nil, err
	}

	parameters.HighLatitudeRule = highLatitudeRule

	// ------------------------------------------------------------
	// Calculate prayer times
	// ------------------------------------------------------------

	prayerTimes, err := calc.NewPrayerTimes(
		coordinates,
		dateComponents,
		parameters,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"calculate prayer times: %w",
			err,
		)
	}

	// ------------------------------------------------------------
	// Convert calculated times to requested timezone.
	// ------------------------------------------------------------

	if err := prayerTimes.SetTimeZone(timezone); err != nil {
		return nil, fmt.Errorf(
			"set prayer timezone: %w",
			err,
		)
	}

	// ------------------------------------------------------------
	// Build our domain model.
	// ------------------------------------------------------------

	schedule := &PrayerSchedule{
		Date:     localDate,
		Timezone: timezone,

		Fajr: PrayerTime{
			Name:       PrayerFajr,
			ArabicName: "الفجر",
			Time:       prayerTimes.Fajr,
		},

		Sunrise: PrayerTime{
			Name:       PrayerSunrise,
			ArabicName: "الشروق",
			Time:       prayerTimes.Sunrise,
		},

		Dhuhr: PrayerTime{
			Name:       PrayerDhuhr,
			ArabicName: "الظهر",
			Time:       prayerTimes.Dhuhr,
		},

		Asr: PrayerTime{
			Name:       PrayerAsr,
			ArabicName: "العصر",
			Time:       prayerTimes.Asr,
		},

		Maghrib: PrayerTime{
			Name:       PrayerMaghrib,
			ArabicName: "المغرب",
			Time:       prayerTimes.Maghrib,
		},

		Isha: PrayerTime{
			Name:       PrayerIsha,
			ArabicName: "العشاء",
			Time:       prayerTimes.Isha,
		},
	}

	return schedule, nil
}

// ------------------------------------------------------------
// Adhan calculation method mapping
// ------------------------------------------------------------

func toAdhanCalculationMethod(
	method CalculationMethod,
) (calc.CalculationMethod, error) {

	switch method {

	case CalculationMethodMuslimWorldLeague:
		return calc.MUSLIM_WORLD_LEAGUE, nil

	case CalculationMethodEgyptian:
		return calc.EGYPTIAN, nil

	case CalculationMethodKarachi:
		return calc.KARACHI, nil

	case CalculationMethodUmmAlQura:
		return calc.UMM_AL_QURA, nil

	case CalculationMethodDubai:
		return calc.DUBAI, nil

	case CalculationMethodMoonSightingCommittee:
		return calc.MOON_SIGHTING_COMMITTEE, nil

	case CalculationMethodNorthAmerica:
		return calc.NORTH_AMERICA, nil

	case CalculationMethodKuwait:
		return calc.KUWAIT, nil

	case CalculationMethodQatar:
		return calc.QATAR, nil

	case CalculationMethodSingapore:
		return calc.SINGAPORE, nil

	default:
		return calc.OTHER, fmt.Errorf(
			"unsupported calculation method %q",
			method,
		)
	}
}

// ------------------------------------------------------------
// Madhab mapping
// ------------------------------------------------------------

func toAdhanMadhab(
	madhab Madhab,
) (calc.AsrJuristicMethod, error) {

	switch madhab {

	case MadhabShafi:
		return calc.SHAFI_HANBALI_MALIKI, nil

	case MadhabHanafi:
		return calc.HANAFI, nil

	default:
		return calc.SHAFI_HANBALI_MALIKI, fmt.Errorf(
			"unsupported madhab %q",
			madhab,
		)
	}
}

// ------------------------------------------------------------
// High latitude mapping
// ------------------------------------------------------------

func toAdhanHighLatitudeRule(
	rule HighLatitudeRule,
) (calc.HighLatitudeRule, error) {

	switch rule {

	case HighLatitudeMiddleOfTheNight:
		return calc.MIDDLE_OF_THE_NIGHT, nil

	case HighLatitudeSeventhOfTheNight:
		return calc.SEVENTH_OF_THE_NIGHT, nil

	case HighLatitudeTwilightAngle:
		return calc.TWILIGHT_ANGLE, nil

	default:
		return calc.NO_HIGH_LATITUDE_RULE, fmt.Errorf(
			"unsupported high latitude rule %q",
			rule,
		)
	}
}
