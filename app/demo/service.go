package demo

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type storager interface {
	Transaction(ctx context.Context, fn func(tx *sql.Tx) error) error
	InsertToDemo(tx *sql.Tx, demo *Demo) error
	InsertToDetail(tx *sql.Tx, detail *Detail) error
}

type Service struct {
	storager storager
}

func NewService(storager storager) *Service {
	return &Service{
		storager: storager,
	}
}

func (s *Service) SaveToDatabase(bufferMessage *[]DemoRequest, syncLock *sync.Mutex) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := new(DemoRequest)
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
		}
		syncLock.Lock()
		*bufferMessage = append(*bufferMessage, DemoRequest{
			Name:    req.Name,
			Surname: req.Surname,
			Age:     req.Age,
			Detail:  req.Detail,
		})
		syncLock.Unlock()
		if bufferMessage != nil && len(*bufferMessage) > 1000 {
			ctxWithOutCancel := context.WithoutCancel(ctx)
			func() {
				err := s.BatchSaveDatabase(ctxWithOutCancel, bufferMessage, syncLock)
				if err != nil {
					// do something when error is appear EX. log to sre, send to kafka to try in case one by one, send back to same kafka
					log.Println(err)
				}
			}()
		}
		// log.Println(len(*bufferMessage))
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Accepted",
		})
	}

}

func (s *Service) BatchSaveDatabase(ctx context.Context, bufferMessage *[]DemoRequest, syncLock *sync.Mutex) error {
	if err := s.storager.Transaction(ctx, func(tx *sql.Tx) error {
		if bufferMessage != nil && len(*bufferMessage) > 0 {
			syncLock.Lock()
			// Unlock when service is end
			defer syncLock.Unlock()

			for _, v := range *bufferMessage {
				demoSave := &Demo{
					Name:    v.Name,
					Surname: v.Surname,
					Age:     v.Age,
				}
				if err := s.storager.InsertToDemo(tx, demoSave); err != nil {
					return err
				}

				detailSave := &Detail{
					DemoID: demoSave.ID,
					Detail: v.Detail,
				}
				if err := s.storager.InsertToDetail(tx, detailSave); err != nil {
					return err
				}
			}
			*bufferMessage = []DemoRequest{}
			log.Println("success Saved")
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
