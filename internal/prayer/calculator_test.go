package prayer

import (
	"testing"
	"time"
)

func TestCalculatorCalculate(t *testing.T) {

	calculator := NewCalculator()

	location := "Asia/Kolkata"

	date := time.Date(
		2026,
		time.October,
		7,
		0,
		0,
		0,
		0,
		time.FixedZone("IST", 5*60*60+30*60),
	)

	settings := CalculationSettings{
		Method:           CalculationMethodKarachi,
		Madhab:           MadhabShafi,
		HighLatitudeRule: HighLatitudeMiddleOfTheNight,
	}

	schedule, err := calculator.Calculate(
		10.5276,
		76.2144,
		date,
		location,
		settings,
	)

	if err != nil {
		t.Fatalf(
			"Calculate() returned error: %v",
			err,
		)
	}

	if schedule == nil {
		t.Fatal("Calculate() returned nil schedule")
	}

	if schedule.Timezone != location {
		t.Fatalf(
			"expected timezone %q, got %q",
			location,
			schedule.Timezone,
		)
	}

	prayers := []PrayerTime{
		schedule.Fajr,
		schedule.Sunrise,
		schedule.Dhuhr,
		schedule.Asr,
		schedule.Maghrib,
		schedule.Isha,
	}

	for _, prayer := range prayers {

		if prayer.Time.IsZero() {
			t.Fatalf(
				"%s returned zero time",
				prayer.Name,
			)
		}

		if prayer.Time.Location().String() != location {
			t.Fatalf(
				"%s has timezone %q, expected %q",
				prayer.Name,
				prayer.Time.Location().String(),
				location,
			)
		}
	}
}

func TestPrayerTimesAreChronological(t *testing.T) {

	calculator := NewCalculator()

	date := time.Date(
		2026,
		time.October,
		7,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	settings := CalculationSettings{
		Method:           CalculationMethodKarachi,
		Madhab:           MadhabShafi,
		HighLatitudeRule: HighLatitudeMiddleOfTheNight,
	}

	schedule, err := calculator.Calculate(
		10.5276,
		76.2144,
		date,
		"Asia/Kolkata",
		settings,
	)

	if err != nil {
		t.Fatalf(
			"Calculate() returned error: %v",
			err,
		)
	}

	times := []struct {
		name string
		time time.Time
	}{
		{"Fajr", schedule.Fajr.Time},
		{"Sunrise", schedule.Sunrise.Time},
		{"Dhuhr", schedule.Dhuhr.Time},
		{"Asr", schedule.Asr.Time},
		{"Maghrib", schedule.Maghrib.Time},
		{"Isha", schedule.Isha.Time},
	}

	for i := 1; i < len(times); i++ {

		previous := times[i-1]
		current := times[i]

		if !current.time.After(previous.time) {
			t.Fatalf(
				"%s (%v) must be after %s (%v)",
				current.name,
				current.time,
				previous.name,
				previous.time,
			)
		}
	}
}

func TestCalculatorRejectsInvalidCoordinates(t *testing.T) {

	calculator := NewCalculator()

	settings := CalculationSettings{
		Method:           CalculationMethodKarachi,
		Madhab:           MadhabShafi,
		HighLatitudeRule: HighLatitudeMiddleOfTheNight,
	}

	_, err := calculator.Calculate(
		100,
		76,
		time.Now(),
		"Asia/Kolkata",
		settings,
	)

	if err == nil {
		t.Fatal("expected invalid latitude error")
	}
}

func TestCalculatorRejectsInvalidTimezone(t *testing.T) {

	calculator := NewCalculator()

	settings := CalculationSettings{
		Method:           CalculationMethodKarachi,
		Madhab:           MadhabShafi,
		HighLatitudeRule: HighLatitudeMiddleOfTheNight,
	}

	_, err := calculator.Calculate(
		10.5276,
		76.2144,
		time.Now(),
		"Invalid/Timezone",
		settings,
	)

	if err == nil {
		t.Fatal("expected invalid timezone error")
	}
}
