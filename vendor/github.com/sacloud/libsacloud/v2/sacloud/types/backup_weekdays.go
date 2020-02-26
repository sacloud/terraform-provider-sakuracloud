// Copyright 2016-2020 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import "sort"

// EBackupSpanWeekday バックアップ取得曜日
type EBackupSpanWeekday string

// BackupSpanWeekdays バックアップ取得曜日
var BackupSpanWeekdays = struct {
	Sunday    EBackupSpanWeekday
	Monday    EBackupSpanWeekday
	Tuesday   EBackupSpanWeekday
	Wednesday EBackupSpanWeekday
	Thursday  EBackupSpanWeekday
	Friday    EBackupSpanWeekday
	Saturday  EBackupSpanWeekday
}{
	Sunday:    EBackupSpanWeekday("sun"),
	Monday:    EBackupSpanWeekday("mon"),
	Tuesday:   EBackupSpanWeekday("tue"),
	Wednesday: EBackupSpanWeekday("wed"),
	Thursday:  EBackupSpanWeekday("thu"),
	Friday:    EBackupSpanWeekday("fri"),
	Saturday:  EBackupSpanWeekday("sat"),
}

// String Stringer実装
func (w EBackupSpanWeekday) String() string {
	return string(w)
}

// BackupSpanWeekdaysOrder バックアップ取得曜日の並び順(日曜開始)
var BackupSpanWeekdaysOrder = map[EBackupSpanWeekday]int{
	BackupSpanWeekdays.Sunday:    0,
	BackupSpanWeekdays.Monday:    1,
	BackupSpanWeekdays.Tuesday:   2,
	BackupSpanWeekdays.Wednesday: 3,
	BackupSpanWeekdays.Thursday:  4,
	BackupSpanWeekdays.Friday:    5,
	BackupSpanWeekdays.Saturday:  6,
}

// SortBackupSpanWeekdays バックアップ取得曜日のソート(日曜開始)
func SortBackupSpanWeekdays(weekdays []EBackupSpanWeekday) {
	sort.Slice(weekdays, func(i, j int) bool {
		return BackupSpanWeekdaysOrder[weekdays[i]] < BackupSpanWeekdaysOrder[weekdays[j]]
	})
}

// BackupWeekdayStrings 有効なバックアップ取得曜日のリスト(文字列)
var BackupWeekdayStrings = []string{
	BackupSpanWeekdays.Sunday.String(),
	BackupSpanWeekdays.Monday.String(),
	BackupSpanWeekdays.Tuesday.String(),
	BackupSpanWeekdays.Wednesday.String(),
	BackupSpanWeekdays.Thursday.String(),
	BackupSpanWeekdays.Friday.String(),
	BackupSpanWeekdays.Saturday.String(),
}
