//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
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
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package models

import (
	"fmt"
	"go.uber.org/atomic"
	"time"
)

type ProcessStats struct {
	SentCount               atomic.Int32
	SentErrorCount          atomic.Int32
	ReceivedExpectedCount   atomic.Int32
	ReceivedUnexpectedCount atomic.Int32
	ReceivedErrorCount      atomic.Int32
}

type ProcessResult struct {
	StartTime *time.Time
	Stats     ProcessStats
	Err       error
}

func (ps ProcessStats) String() string {
	return fmt.Sprintf("sent messages: %d | send errors: %d | expected received messages: %d | unexpected received messages: %d | receive errors: %d",
		ps.SentCount, ps.SentErrorCount, ps.ReceivedExpectedCount, ps.ReceivedUnexpectedCount, ps.ReceivedErrorCount)
}
