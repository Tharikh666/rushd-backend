package prayer

import (
	"context"
	"fmt"
	"time"

	prayerv1 "rushd-backend/gen/prayer/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements the RUSHD PrayerService gRPC API.
type Service struct {
	prayerv1.UnimplementedPrayerServiceServer

	calculator *Calculator
}

// NewService creates a new PrayerService.
func NewService(calculator *Calculator) *Service {
	return &Service{
		calculator: calculator,
	}
}

// GetPrayerTimes calculates and returns prayer times for the requested
// location, date, timezone, and calculation settings.
func (s *Service) GetPrayerTimes(
	ctx context.Context,
	req *prayerv1.GetPrayerTimesRequest,
) (*prayerv1.GetPrayerTimesResponse, error) {

	// Context cancellation/deadline check.
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, ctx.Err().Error())
	default:
	}

	// Validate request.
	if req == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"request is required",
		)
	}

	if req.Calculation == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"calculation settings are required",
		)
	}

	// Validate timezone.
	if req.Timezone == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"timezone is required",
		)
	}

	if _, err := time.LoadLocation(req.Timezone); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid timezone %q: %v",
			req.Timezone,
			err,
		)
	}

	// Validate date.
	if req.Date == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"date is required",
		)
	}

	location, err := time.LoadLocation(req.Timezone)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid timezone %q: %v",
			req.Timezone,
			err,
		)
	}

	date, err := time.ParseInLocation(
		"2006-01-02",
		req.Date,
		location,
	)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid date %q: expected YYYY-MM-DD",
			req.Date,
		)
	}

	// Convert protobuf calculation settings into domain settings.
	settings, err := calculationSettingsFromProto(req.Calculation)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	// Make sure the calculator is configured.
	if s.calculator == nil {
		return nil, status.Error(
			codes.Internal,
			"prayer calculator is not configured",
		)
	}

	// Calculate prayer times.
	schedule, err := s.calculator.Calculate(
		req.Latitude,
		req.Longitude,
		date,
		req.Timezone,
		settings,
	)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	// Convert domain response into protobuf response.
	return prayerScheduleToProto(schedule), nil
}

// calculationSettingsFromProto converts protobuf calculation settings
// into the internal domain model.
func calculationSettingsFromProto(
	settings *prayerv1.CalculationSettings,
) (CalculationSettings, error) {

	if settings == nil {
		return CalculationSettings{}, fmt.Errorf(
			"calculation settings are required",
		)
	}

	method, err := calculationMethodFromProto(settings.Method)
	if err != nil {
		return CalculationSettings{}, err
	}

	madhab, err := madhabFromProto(settings.Madhab)
	if err != nil {
		return CalculationSettings{}, err
	}

	highLatitudeRule, err := highLatitudeRuleFromProto(
		settings.HighLatitudeRule,
	)
	if err != nil {
		return CalculationSettings{}, err
	}

	return CalculationSettings{
		Method:           method,
		Madhab:           madhab,
		HighLatitudeRule: highLatitudeRule,
	}, nil
}

// calculationMethodFromProto converts the protobuf calculation method
// into the internal domain calculation method.
func calculationMethodFromProto(
	method prayerv1.CalculationMethod,
) (CalculationMethod, error) {

	switch method {

	case prayerv1.CalculationMethod_MUSLIM_WORLD_LEAGUE:
		return CalculationMethodMuslimWorldLeague, nil

	case prayerv1.CalculationMethod_EGYPTIAN:
		return CalculationMethodEgyptian, nil

	case prayerv1.CalculationMethod_KARACHI:
		return CalculationMethodKarachi, nil

	case prayerv1.CalculationMethod_UMM_AL_QURA:
		return CalculationMethodUmmAlQura, nil

	case prayerv1.CalculationMethod_DUBAI:
		return CalculationMethodDubai, nil

	case prayerv1.CalculationMethod_MOON_SIGHTING_COMMITTEE:
		return CalculationMethodMoonSightingCommittee, nil

	case prayerv1.CalculationMethod_NORTH_AMERICA:
		return CalculationMethodNorthAmerica, nil

	case prayerv1.CalculationMethod_KUWAIT:
		return CalculationMethodKuwait, nil

	case prayerv1.CalculationMethod_QATAR:
		return CalculationMethodQatar, nil

	case prayerv1.CalculationMethod_SINGAPORE:
		return CalculationMethodSingapore, nil

	default:
		return "", fmt.Errorf(
			"unsupported calculation method %q",
			method.String(),
		)
	}
}

// madhabFromProto converts the protobuf Madhab enum into the internal
// domain Madhab.
func madhabFromProto(
	madhab prayerv1.Madhab,
) (Madhab, error) {

	switch madhab {

	case prayerv1.Madhab_SHAFI:
		return MadhabShafi, nil

	case prayerv1.Madhab_HANAFI:
		return MadhabHanafi, nil

	default:
		return "", fmt.Errorf(
			"unsupported madhab %q",
			madhab.String(),
		)
	}
}

// highLatitudeRuleFromProto converts the protobuf high-latitude rule
// into the internal domain high-latitude rule.
func highLatitudeRuleFromProto(
	rule prayerv1.HighLatitudeRule,
) (HighLatitudeRule, error) {

	switch rule {

	case prayerv1.HighLatitudeRule_MIDDLE_OF_THE_NIGHT:
		return HighLatitudeMiddleOfTheNight, nil

	case prayerv1.HighLatitudeRule_SEVENTH_OF_THE_NIGHT:
		return HighLatitudeSeventhOfTheNight, nil

	case prayerv1.HighLatitudeRule_TWILIGHT_ANGLE:
		return HighLatitudeTwilightAngle, nil

	default:
		return "", fmt.Errorf(
			"unsupported high latitude rule %q",
			rule.String(),
		)
	}
}

// prayerScheduleToProto converts the internal prayer schedule into
// the protobuf response.
func prayerScheduleToProto(
	schedule *PrayerSchedule,
) *prayerv1.GetPrayerTimesResponse {

	return &prayerv1.GetPrayerTimesResponse{
		Date:     schedule.Date.Format("2006-01-02"),
		Timezone: schedule.Timezone,

		Fajr: prayerTimeToProto(schedule.Fajr),

		Sunrise: prayerTimeToProto(schedule.Sunrise),

		Dhuhr: prayerTimeToProto(schedule.Dhuhr),

		Asr: prayerTimeToProto(schedule.Asr),

		Maghrib: prayerTimeToProto(schedule.Maghrib),

		Isha: prayerTimeToProto(schedule.Isha),
	}
}

// prayerTimeToProto converts an internal PrayerTime into the
// corresponding protobuf PrayerTime.
func prayerTimeToProto(
	prayerTime PrayerTime,
) *prayerv1.PrayerTime {

	return &prayerv1.PrayerTime{
		Name:        string(prayerTime.Name),
		ArabicName:  prayerTime.ArabicName,
		Time:        prayerTime.Time.Format("15:04"),
		IsoDatetime: prayerTime.Time.Format(time.RFC3339),
	}
}
