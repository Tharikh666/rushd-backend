package prayer

import "time"

// CalculationSettings controls how prayer times are calculated.
type CalculationSettings struct {
	Method           CalculationMethod
	Madhab           Madhab
	HighLatitudeRule HighLatitudeRule
}

// CalculationMethod represents an Islamic prayer calculation convention.
type CalculationMethod string

const (
	CalculationMethodMuslimWorldLeague     CalculationMethod = "MUSLIM_WORLD_LEAGUE"
	CalculationMethodEgyptian              CalculationMethod = "EGYPTIAN"
	CalculationMethodKarachi               CalculationMethod = "KARACHI"
	CalculationMethodUmmAlQura             CalculationMethod = "UMM_AL_QURA"
	CalculationMethodDubai                 CalculationMethod = "DUBAI"
	CalculationMethodMoonSightingCommittee CalculationMethod = "MOON_SIGHTING_COMMITTEE"
	CalculationMethodNorthAmerica          CalculationMethod = "NORTH_AMERICA"
	CalculationMethodKuwait                CalculationMethod = "KUWAIT"
	CalculationMethodQatar                 CalculationMethod = "QATAR"
	CalculationMethodSingapore             CalculationMethod = "SINGAPORE"
)

// Madhab controls the juristic method used for Asr.
type Madhab string

const (
	MadhabShafi  Madhab = "SHAFI"
	MadhabHanafi Madhab = "HANAFI"
)

// HighLatitudeRule controls how Fajr and Isha are bounded
// at high latitudes.
type HighLatitudeRule string

const (
	HighLatitudeMiddleOfTheNight  HighLatitudeRule = "MIDDLE_OF_THE_NIGHT"
	HighLatitudeSeventhOfTheNight HighLatitudeRule = "SEVENTH_OF_THE_NIGHT"
	HighLatitudeTwilightAngle     HighLatitudeRule = "TWILIGHT_ANGLE"
)

// PrayerName identifies a prayer event.
type PrayerName string

const (
	PrayerFajr    PrayerName = "Fajr"
	PrayerSunrise PrayerName = "Sunrise"
	PrayerDhuhr   PrayerName = "Dhuhr"
	PrayerAsr     PrayerName = "Asr"
	PrayerMaghrib PrayerName = "Maghrib"
	PrayerIsha    PrayerName = "Isha"
)

// PrayerTime represents one calculated prayer event.
type PrayerTime struct {
	Name       PrayerName
	ArabicName string
	Time       time.Time
}

// PrayerSchedule contains all calculated prayer events for one day.
type PrayerSchedule struct {
	Date     time.Time
	Timezone string

	Fajr    PrayerTime
	Sunrise PrayerTime
	Dhuhr   PrayerTime
	Asr     PrayerTime
	Maghrib PrayerTime
	Isha    PrayerTime
}
