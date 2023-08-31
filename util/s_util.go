package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/constraints"

	sync "github.com/sasha-s/go-deadlock"
)

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func Truncf(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func MinInt(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func AbsInt(i int) int {
	if i > 0 {
		return i
	}
	return -i
}

// jsonifyWhatever converts whatever is passed into a JSON string.
func JsonifyWhatever(i interface{}) string {
	jsonb, err := json.Marshal(i)
	if err != nil {
		log.Panic(err)
	}
	return string(jsonb)
}

// jsonifyWhateverToBytes converts whatever is passed into a JSON byte slice.
func JsonifyWhateverToBytes(i interface{}) []byte {
	jsonb, err := json.Marshal(i)
	if err != nil {
		log.Panic(err)
	}
	return jsonb
}

// jsonifyWhateverToBuffer converts whatever is passed into a
// JSON byte buffer.
func JsonifyWhateverToBuffer(i interface{}) *bytes.Buffer {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(i)
	return b
}

// WithMutex extends the Mutex type with the convenient .With(func) function
type WithMutex struct {
	sync.Mutex
}

// WithLock executes the given function with the mutex locked
func (m *WithMutex) WithLock(f func()) {
	m.Mutex.Lock()
	f()
	m.Mutex.Unlock()
}

// WithRWMutex extends the RWMutex type with convenient .With(func) functions
type WithRWMutex struct {
	sync.RWMutex
}

// WithRLock executes the given function with the mutex rlocked
func (m *WithRWMutex) WithRLock(f func()) {
	m.RWMutex.RLock()
	f()
	m.RWMutex.RUnlock()
}

// WithWLock executes the given function with the mutex wlocked
func (m *WithRWMutex) WithWLock(f func()) {
	m.RWMutex.Lock()
	f()
	m.RWMutex.Unlock()
}

// Converts the given Unix timestamp to time.Time
func UnixTimeStampToUTCTime(ts int) time.Time {
	return time.Unix(int64(ts), 0)
}

// Gets the current Unix timestamp in UTC
func GetNowUTC() int64 {
	return time.Now().UTC().Unix()
}

// Mashals the given map of strings to JSON
func stringMap2JsonBytes(m map[string]string) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Panicln("Cannot json-ise the map:", err)
	}
	return b
}

// Returns a hex-encoded hash of the given byte slice
func HashBytesToHexString(b []byte) string {
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:])
}

// Returns a hex-encoded hash of the given file
func HashFileToHexString(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func NowUTC() time.Time {
	return time.Now().UTC()
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func RoundF32toInt(f float32) int {
	return int(math.Round(float64(f)))
}

func euclidDistance(lat1, lng1, lat2, lng2 float32) float32 {
	dLat := float64(lat2 - lat1)
	dLng := float64(lng2 - lng1)
	return float32(math.Sqrt(dLat*dLat + dLng*dLng))
}

func IfToFloat64(i interface{}) (f float64) {
	if i != nil {
		switch v := i.(type) {
		case float64:
			f = v
		case int32:
			f = float64(v)
		case int64:
			f = float64(v)
		case int:
			f = float64(v)
		}
	}
	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.!")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetHTTPJSONdict(url string) (m map[string]interface{}, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error retrieving JSON document at %s: %s", url, res.Status)
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	m = map[string]interface{}{}
	err = json.Unmarshal(data, &m)
	return
}

func GetHTTPJSON(url string, i interface{}) (err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error retrieving JSON document at %s: %s", url, res.Status)
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, i)
	return
}

type Number interface {
	constraints.Integer | constraints.Float
}

func Abs[T Number](n T) T {
	if n < T(0) {
		return -n
	} else {
		return n
	}
}

type BiggishNumber interface {
	~uint | ~uint32 | ~uint64 | ~uintptr | constraints.Float
}

func BToMB[T BiggishNumber](n T) T {
	return n / T(1024*1024)
}
