package prayer

import (
	"context"
	"testing"
	"time"

	prayerv1 "rushd-backend/gen/prayer/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServiceGetPrayerTimes(t *testing.T) {
	service := NewService(NewCalculator())

	req := &prayerv1.GetPrayerTimesRequest{
		Latitude:  10.5276,
		Longitude: 76.2144,
		Timezone:  "Asia/Kolkata",
		Date:      "2026-10-07",
		Calculation: &prayerv1.CalculationSettings{
			Method:           prayerv1.CalculationMethod_KARACHI,
			Madhab:           prayerv1.Madhab_SHAFI,
			HighLatitudeRule: prayerv1.HighLatitudeRule_MIDDLE_OF_THE_NIGHT,
		},
	}

	response, err := service.GetPrayerTimes(context.Background(), req)
	if err != nil {
		t.Fatalf("GetPrayerTimes() returned error: %v", err)
	}

	if response == nil {
		t.Fatal("GetPrayerTimes() returned nil response")
	}

	if response.Date != "2026-10-07" {
		t.Errorf(
			"expected date %q, got %q",
			"2026-10-07",
			response.Date,
		)
	}

	if response.Timezone != "Asia/Kolkata" {
		t.Errorf(
			"expected timezone %q, got %q",
			"Asia/Kolkata",
			response.Timezone,
		)
	}

	prayers := []*prayerv1.PrayerTime{
		response.Fajr,
		response.Sunrise,
		response.Dhuhr,
		response.Asr,
		response.Maghrib,
		response.Isha,
	}

	for _, prayer := range prayers {
		if prayer == nil {
			t.Fatal("response contains nil prayer time")
		}

		if prayer.Name == "" {
			t.Error("prayer name is empty")
		}

		if prayer.Time == "" {
			t.Error("prayer time is empty")
		}

		if prayer.IsoDatetime == "" {
			t.Error("prayer ISO datetime is empty")
		}
	}
}

func TestServiceGetPrayerTimesRejectsNilRequest(t *testing.T) {
	service := NewService(NewCalculator())

	_, err := service.GetPrayerTimes(
		context.Background(),
		nil,
	)

	if err == nil {
		t.Fatal("expected error for nil request")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf(
			"expected InvalidArgument, got %s",
			status.Code(err),
		)
	}
}

func TestServiceGetPrayerTimesRejectsMissingCalculation(
	t *testing.T,
) {
	service := NewService(NewCalculator())

	req := &prayerv1.GetPrayerTimesRequest{
		Latitude:  10.5276,
		Longitude: 76.2144,
		Timezone:  "Asia/Kolkata",
		Date:      "2026-10-07",
	}

	_, err := service.GetPrayerTimes(
		context.Background(),
		req,
	)

	if err == nil {
		t.Fatal("expected error for missing calculation settings")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf(
			"expected InvalidArgument, got %s",
			status.Code(err),
		)
	}
}

func TestServiceGetPrayerTimesRejectsInvalidTimezone(
	t *testing.T,
) {
	service := NewService(NewCalculator())

	req := &prayerv1.GetPrayerTimesRequest{
		Latitude:  10.5276,
		Longitude: 76.2144,
		Timezone:  "Invalid/Timezone",
		Date:      "2026-10-07",
		Calculation: &prayerv1.CalculationSettings{
			Method:           prayerv1.CalculationMethod_KARACHI,
			Madhab:           prayerv1.Madhab_SHAFI,
			HighLatitudeRule: prayerv1.HighLatitudeRule_MIDDLE_OF_THE_NIGHT,
		},
	}

	_, err := service.GetPrayerTimes(
		context.Background(),
		req,
	)

	if err == nil {
		t.Fatal("expected error for invalid timezone")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf(
			"expected InvalidArgument, got %s",
			status.Code(err),
		)
	}
}

func TestServiceGetPrayerTimesRejectsInvalidDate(
	t *testing.T,
) {
	service := NewService(NewCalculator())

	req := &prayerv1.GetPrayerTimesRequest{
		Latitude:  10.5276,
		Longitude: 76.2144,
		Timezone:  "Asia/Kolkata",
		Date:      "not-a-date",
		Calculation: &prayerv1.CalculationSettings{
			Method:           prayerv1.CalculationMethod_KARACHI,
			Madhab:           prayerv1.Madhab_SHAFI,
			HighLatitudeRule: prayerv1.HighLatitudeRule_MIDDLE_OF_THE_NIGHT,
		},
	}

	_, err := service.GetPrayerTimes(
		context.Background(),
		req,
	)

	if err == nil {
		t.Fatal("expected error for invalid date")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf(
			"expected InvalidArgument, got %s",
			status.Code(err),
		)
	}
}

func TestPrayerResponseTimesAreChronological(t *testing.T) {
	service := NewService(NewCalculator())

	req := &prayerv1.GetPrayerTimesRequest{
		Latitude:  10.5276,
		Longitude: 76.2144,
		Timezone:  "Asia/Kolkata",
		Date:      "2026-10-07",
		Calculation: &prayerv1.CalculationSettings{
			Method:           prayerv1.CalculationMethod_KARACHI,
			Madhab:           prayerv1.Madhab_SHAFI,
			HighLatitudeRule: prayerv1.HighLatitudeRule_MIDDLE_OF_THE_NIGHT,
		},
	}

	response, err := service.GetPrayerTimes(
		context.Background(),
		req,
	)

	if err != nil {
		t.Fatalf("GetPrayerTimes() returned error: %v", err)
	}

	parse := func(value string) time.Time {
		result, err := time.Parse(
			"15:04",
			value,
		)

		if err != nil {
			t.Fatalf(
				"failed to parse prayer time %q: %v",
				value,
				err,
			)
		}

		return result
	}

	times := []struct {
		name string
		time string
	}{
		{"Fajr", response.Fajr.Time},
		{"Sunrise", response.Sunrise.Time},
		{"Dhuhr", response.Dhuhr.Time},
		{"Asr", response.Asr.Time},
		{"Maghrib", response.Maghrib.Time},
		{"Isha", response.Isha.Time},
	}

	for i := 1; i < len(times); i++ {
		previous := parse(times[i-1].time)
		current := parse(times[i].time)

		if !current.After(previous) {
			t.Errorf(
				"%s (%s) should be after %s (%s)",
				times[i].name,
				times[i].time,
				times[i-1].name,
				times[i-1].time,
			)
		}
	}
}
