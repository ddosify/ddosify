/*
*
*	Ddosify - Load testing tool for any web system.
*   Copyright (C) 2021  Ddosify (https://ddosify.com)
*
*   This program is free software: you can redistribute it and/or modify
*   it under the terms of the GNU Affero General Public License as published
*   by the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.
*
*   This program is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU Affero General Public License for more details.
*
*   You should have received a copy of the GNU Affero General Public License
*   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*
 */

package core

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"

	"ddosify.com/hammer/core/proxy"
	"ddosify.com/hammer/core/report"
	"ddosify.com/hammer/core/scenario"
	"ddosify.com/hammer/core/types"
)

const (
	// interval in milisecond
	tickerInterval = 100
)

type engine struct {
	hammer types.Hammer

	proxyService    proxy.ProxyService
	scenarioService *scenario.ScenarioService
	reportService   report.ReportService

	tickCounter int
	reqCountArr []int
	wg          sync.WaitGroup

	responseChan chan *types.Response

	ctx context.Context
}

func NewEngine(ctx context.Context, h types.Hammer) (e *engine, err error) {
	err = h.Validate()
	if err != nil {
		return
	}

	ps, err := proxy.NewProxyService(h.Proxy.Strategy)
	if err != nil {
		return
	}

	rs, err := report.NewReportService(h.ReportDestination)
	if err != nil {
		return
	}

	ss := scenario.NewScenarioService()

	e = &engine{
		hammer:          h,
		ctx:             ctx,
		proxyService:    ps,
		scenarioService: ss,
		reportService:   rs,
	}

	return
}

func (e *engine) Init() (err error) {
	if err = e.proxyService.Init(e.hammer.Proxy); err != nil {
		return
	}

	if err = e.scenarioService.Init(e.hammer.Scenario, e.proxyService.GetAll(), e.ctx); err != nil {
		return
	}

	if err = e.reportService.Init(); err != nil {
		return
	}

	e.initReqCountArr()
	return
}

func (e *engine) Start() {
	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Millisecond)
	e.responseChan = make(chan *types.Response, e.hammer.TotalReqCount)
	go e.reportService.Start(e.responseChan)

	defer func() {
		ticker.Stop()
		e.stop()
	}()

	e.tickCounter = 0
	e.wg = sync.WaitGroup{}
	var mutex = &sync.Mutex{}
	for range ticker.C {
		if e.tickCounter >= len(e.reqCountArr) {
			return
		}

		select {
		case <-e.ctx.Done():
			return
		default:
			mutex.Lock()
			e.wg.Add(e.reqCountArr[e.tickCounter])
			go e.runWorkers(e.tickCounter)
			e.tickCounter++
			mutex.Unlock()
		}
	}
}

func (e *engine) runWorkers(c int) {
	for i := 1; i <= e.reqCountArr[c]; i++ {
		go func() {
			e.runWorker()
			e.wg.Done()
		}()
	}
}

func (e *engine) runWorker() {
	p := e.proxyService.GetProxy()
	res, err := e.scenarioService.Do(p)

	if err != nil && err.Type == types.ErrorProxy {
		e.proxyService.ReportProxy(p, err.Reason)
	}
	if err != nil && err.Type == types.ErrorIntented {
		// Don't report intentionally created errors. Like canceled requests.
		return
	}

	e.responseChan <- res
}

func (e *engine) stop() {
	e.wg.Wait()
	close(e.responseChan)
	<-e.reportService.DoneChan()
	e.reportService.Report()
}

func (e *engine) initReqCountArr() {
	if e.hammer.TimeReqCountMap != nil {
		fmt.Println("initReqCountArr from TimeReqCountMap")
	} else {
		length := int(e.hammer.TestDuration * int(time.Second/(tickerInterval*time.Millisecond)))
		e.reqCountArr = make([]int, length)

		switch e.hammer.LoadType {
		case types.LoadTypeLinear:
			e.createLinearReqCountArr()
		case types.LoadTypeIncremental:
			e.createIncrementalReqCountArr()
		case types.LoadTypeWaved:
			e.createWavedReqCountArr()
		}
	}
}

func (e *engine) createLinearReqCountArr() {
	createLinearDistArr(e.hammer.TotalReqCount, e.reqCountArr)
}

func (e *engine) createIncrementalReqCountArr() {
	steps := createIncrementalDistArr(e.hammer.TotalReqCount, e.hammer.TestDuration)
	tickPerSecond := int(time.Second / (tickerInterval * time.Millisecond))
	for i := range steps {
		tickArrStartIndex := i * tickPerSecond
		tickArrEndIndex := tickArrStartIndex + tickPerSecond
		segment := e.reqCountArr[tickArrStartIndex:tickArrEndIndex]
		createLinearDistArr(steps[i], segment)
	}
}

func (e *engine) createWavedReqCountArr() {
	tickPerSecond := int(time.Second / (tickerInterval * time.Millisecond))
	quarterWaveCount := int((math.Log2(float64(e.hammer.TestDuration))))
	qWaveDuration := int(e.hammer.TestDuration / quarterWaveCount)
	reqCountPerQWave := int(e.hammer.TotalReqCount / quarterWaveCount)
	tickArrStartIndex := 0

	for i := 0; i < quarterWaveCount; i++ {
		if i == quarterWaveCount-1 {
			// Add remaining req count to the last wave
			reqCountPerQWave += e.hammer.TotalReqCount - (reqCountPerQWave * quarterWaveCount)
		}

		steps := createIncrementalDistArr(reqCountPerQWave, qWaveDuration)
		if i%2 == 1 {
			reverse(steps)
		}

		for j := range steps {
			tickArrEndIndex := tickArrStartIndex + tickPerSecond
			segment := e.reqCountArr[tickArrStartIndex:tickArrEndIndex]
			createLinearDistArr(steps[j], segment)
			tickArrStartIndex += tickPerSecond
		}
	}
}

func createLinearDistArr(count int, arr []int) {
	len := len(arr)
	minReqCount := int(count / len)
	remaining := count - minReqCount*len
	for i := range arr {
		plusOne := 0
		if i < remaining {
			plusOne = 1
		}
		reqCount := minReqCount + plusOne
		arr[i] = reqCount
	}
}

func createIncrementalDistArr(count int, len int) []int {
	steps := make([]int, len)
	sum := (len * (len + 1)) / 2
	incrementStep := int(math.Ceil(float64(sum) / float64(count)))
	val := 0
	for i := range steps {
		if i > 0 {
			val = steps[i-1]
		}

		if i%incrementStep == 0 {
			steps[i] = val + 1
		} else {
			steps[i] = val
		}
	}

	sum = arraySum(steps)

	factor := count / sum
	remaining := count - (sum * factor)
	plus := remaining / len
	lastRemaining := remaining - (plus * len)
	for i := range steps {
		steps[i] = steps[i]*factor + plus
		if len-i-1 < lastRemaining {
			steps[i]++
		}
	}
	return steps
}

func arraySum(steps []int) int {
	sum := 0
	for i := range steps {
		sum += steps[i]
	}
	return sum
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
