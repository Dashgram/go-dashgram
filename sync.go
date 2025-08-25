package dashgram

import "context"

func (d *Dashgram) TrackEventWithContext(ctx context.Context, event any) error {
	if d.useAsync {
		d.TrackEventAsyncWithContext(ctx, event)
		return nil
	}

	requestData := TrackEventRequest{
		Origin:  d.Origin,
		Updates: []any{event},
	}

	return d.request(ctx, "track", requestData)
}

func (d *Dashgram) InvitedByWithContext(ctx context.Context, userID int, invitedBy int) error {
	if d.useAsync {
		d.InvitedByAsyncWithContext(ctx, userID, invitedBy)
		return nil
	}

	requestData := InvitedByRequest{
		UserID:    userID,
		InvitedBy: invitedBy,
		Origin:    d.Origin,
	}

	return d.request(ctx, "invited_by", requestData)
}

func (d *Dashgram) TrackEvent(event any) error {
	return d.TrackEventWithContext(context.Background(), event)
}

func (d *Dashgram) InvitedBy(userID int, invitedBy int) error {
	return d.InvitedByWithContext(context.Background(), userID, invitedBy)
}
