package dashgram

import "context"

func (d *Dashgram) enqueueTask(task asyncTask) {
	select {
	case d.taskChan <- task:
		// Task enqueued successfully
	case <-d.workerCtx.Done():
		// Worker is shutting down, task dropped
	}
}

// TrackEventAsync enqueues an event tracking task to be processed asynchronously
func (d *Dashgram) TrackEventAsyncWithContext(ctx context.Context, event any) {
	requestData := TrackEventRequest{
		Origin:  d.Origin,
		Updates: []any{event},
	}

	d.enqueueTask(asyncTask{
		ctx:      ctx,
		endpoint: "track",
		data:     requestData,
	})
}

// InvitedByAsync enqueues an invitation tracking task to be processed asynchronously
func (d *Dashgram) InvitedByAsyncWithContext(ctx context.Context, userID int, invitedBy int) {
	requestData := InvitedByRequest{
		UserID:    userID,
		InvitedBy: invitedBy,
		Origin:    d.Origin,
	}

	d.enqueueTask(asyncTask{
		ctx:      ctx,
		endpoint: "invited_by",
		data:     requestData,
	})
}

func (d *Dashgram) TrackEventAsync(event any) {
	d.TrackEventAsyncWithContext(context.Background(), event)
}

func (d *Dashgram) InvitedByAsync(userID int, invitedBy int) {
	d.InvitedByAsyncWithContext(context.Background(), userID, invitedBy)
}
