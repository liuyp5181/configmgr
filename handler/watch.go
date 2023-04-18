package handler

import (
	"github.com/liuyp5181/base/log"
	pb "github.com/liuyp5181/configmgr/api"
	"sync"
)

var watcher = struct {
	mux sync.RWMutex
	chs map[string][]chan *pb.WatchRes
}{chs: map[string][]chan *pb.WatchRes{}}

func PutWatch(key string, val []byte) {
	watcher.mux.RLock()
	defer watcher.mux.RUnlock()
	for _, ch := range watcher.chs[key] {
		ch <- &pb.WatchRes{
			Type: pb.WatchType_PUT,
			Key:  key,
			Val:  val,
		}
	}
}

func DeleteWatch(key string) {
	watcher.mux.RLock()
	defer watcher.mux.RUnlock()
	for _, ch := range watcher.chs[key] {
		ch <- &pb.WatchRes{
			Type: pb.WatchType_DELETE,
			Key:  key,
		}
	}
}

func (s *ServiceImpl) Watch(req *pb.WatchReq, stream pb.Greeter_WatchServer) (err error) {
	ch := make(chan *pb.WatchRes, 10)
	watcher.mux.Lock()
	watcher.chs[req.Key] = append(watcher.chs[req.Key], ch)
	watcher.mux.Unlock()
	for {
		err = stream.Send(<-ch)
		if err != nil {
			log.Error("Send err = ", err)
			close(ch)
			watcher.mux.Lock()
			for i, c := range watcher.chs[req.Key] {
				if c == ch {
					watcher.chs[req.Key] = append(watcher.chs[req.Key][:i], watcher.chs[req.Key][i+1:]...)
					break
				}
			}
			watcher.mux.Unlock()
			return err
		}
	}

}
