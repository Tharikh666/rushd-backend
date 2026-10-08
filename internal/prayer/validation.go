package prayer

import (
	"fmt"
	"time"
)

func validateCoordinates(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf(
			"latitude must be between -90 and 90",
		)
	}

	if longitude < -180 || longitude > 180 {
		return fmt.Errorf(
			"longitude must be between -180 and 180",
		)
	}

	return nil
}

func validateTimezone(timezone string) (*time.Location, error) {
	if timezone == "" {
		return nil, fmt.Errorf("timezone is required")
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid timezone %q: %w",
			timezone,
			err,
		)
	}

	return location, nil
}

func validateDate(date time.Time) error {
	if date.IsZero() {
		return fmt.Errorf("date is required")
	}

	return nil
}

func validateCalculationSettings(
	settings CalculationSettings,
) error {

	switch settings.Method {
	case CalculationMethodMuslimWorldLeague,
		CalculationMethodEgyptian,
		CalculationMethodKarachi,
		CalculationMethodUmmAlQura,
		CalculationMethodDubai,
		CalculationMethodMoonSightingCommittee,
		CalculationMethodNorthAmerica,
		CalculationMethodKuwait,
		CalculationMethodQatar,
		CalculationMethodSingapore:
	default:
		return fmt.Errorf(
			"unsupported calculation method %q",
			settings.Method,
		)
	}

	switch settings.Madhab {
	case MadhabShafi, MadhabHanafi:
	default:
		return fmt.Errorf(
			"unsupported madhab %q",
			settings.Madhab,
		)
	}

	switch settings.HighLatitudeRule {
	case HighLatitudeMiddleOfTheNight,
		HighLatitudeSeventhOfTheNight,
		HighLatitudeTwilightAngle:
	default:
		return fmt.Errorf(
			"unsupported high latitude rule %q",
			settings.HighLatitudeRule,
		)
	}

	return nil
}
