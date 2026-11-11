package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

type ConversationRepository struct{ db *sql.DB }

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) List(ctx context.Context, tenantID uint64, status string, limit, offset int) ([]domain.Conversation, int, error) {
	where := "c.tenant_id = ?"
	args := []any{tenantID}
	if status != "" {
		where += " AND c.status = ?"
		args = append(args, status)
	}

	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM conversations c WHERE "+where, args...).Scan(&total)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id,c.tenant_id,c.inbox_id,c.contact_id,c.assigned_user_id,c.status,
		       c.unread_count,c.last_message_at,c.created_at,c.updated_at,
		       ct.name,ct.phone,i.name
		FROM conversations c
		JOIN contacts ct ON ct.id=c.contact_id
		JOIN inboxes i ON i.id=c.inbox_id
		WHERE `+where+` ORDER BY c.last_message_at DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("conversations.List: %w", err)
	}
	defer rows.Close()

	var convs []domain.Conversation
	for rows.Next() {
		var cv domain.Conversation
		var contactName, contactPhone, inboxName string
		err = rows.Scan(&cv.ID, &cv.TenantID, &cv.InboxID, &cv.ContactID,
			&cv.AssignedUserID, &cv.Status, &cv.UnreadCount, &cv.LastMessageAt,
			&cv.CreatedAt, &cv.UpdatedAt,
			&contactName, &contactPhone, &inboxName)
		if err != nil {
			return nil, 0, err
		}
		cv.Contact = &domain.Contact{Name: contactName, Phone: contactPhone}
		cv.Inbox = &domain.Inbox{Name: inboxName}
		if cv.LastMessageAt.Valid {
			cv.LastMessageAtV = &cv.LastMessageAt.Time
		}
		convs = append(convs, cv)
	}
	return convs, total, nil
}

func (r *ConversationRepository) FindByID(ctx context.Context, id, tenantID uint64) (*domain.Conversation, error) {
	var cv domain.Conversation
	err := r.db.QueryRowContext(ctx,
		"SELECT id,tenant_id,inbox_id,contact_id,assigned_user_id,status,unread_count,last_message_at,created_at,updated_at FROM conversations WHERE id=? AND tenant_id=?",
		id, tenantID).Scan(&cv.ID, &cv.TenantID, &cv.InboxID, &cv.ContactID,
		&cv.AssignedUserID, &cv.Status, &cv.UnreadCount, &cv.LastMessageAt,
		&cv.CreatedAt, &cv.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &cv, err
}

func (r *ConversationRepository) Create(ctx context.Context, cv domain.Conversation) (*domain.Conversation, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO conversations(tenant_id,inbox_id,contact_id,status) VALUES(?,?,?,'open')",
		cv.TenantID, cv.InboxID, cv.ContactID)
	if err != nil {
		return nil, fmt.Errorf("conversations.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return r.FindByID(ctx, uint64(id), cv.TenantID)
}

func (r *ConversationRepository) UpdateStatus(ctx context.Context, id, tenantID uint64, status string, assignedUserID *uint64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE conversations SET status=?,assigned_user_id=?,updated_at=NOW() WHERE id=? AND tenant_id=?",
		status, assignedUserID, id, tenantID)
	return err
}

func (r *ConversationRepository) ListMessages(ctx context.Context, conversationID, tenantID uint64, limit, offset int) ([]domain.Message, int, error) {
	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM messages WHERE conversation_id=? AND tenant_id=?", conversationID, tenantID).Scan(&total)

	rows, err := r.db.QueryContext(ctx,
		"SELECT id,conversation_id,tenant_id,sender_type,sender_id,content,content_type,COALESCE(media_url,''),is_read,created_at FROM messages WHERE conversation_id=? AND tenant_id=? ORDER BY created_at ASC LIMIT ? OFFSET ?",
		conversationID, tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("messages.List: %w", err)
	}
	defer rows.Close()

	var msgs []domain.Message
	for rows.Next() {
		var m domain.Message
		var senderID sql.NullInt64
		if err = rows.Scan(&m.ID, &m.ConversationID, &m.TenantID, &m.SenderType, &senderID,
			&m.Content, &m.ContentType, &m.MediaURL, &m.IsRead, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		if senderID.Valid {
			v := uint64(senderID.Int64)
			m.SenderID = &v
		}
		msgs = append(msgs, m)
	}
	return msgs, total, nil
}

func (r *ConversationRepository) CreateMessage(ctx context.Context, m domain.Message) (*domain.Message, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO messages(conversation_id,tenant_id,sender_type,sender_id,content,content_type) VALUES(?,?,?,?,?,?)",
		m.ConversationID, m.TenantID, m.SenderType, m.SenderID, m.Content, m.ContentType)
	if err != nil {
		return nil, fmt.Errorf("messages.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	m.ID = uint64(id)

	// Update conversation last_message_at
	r.db.ExecContext(ctx, "UPDATE conversations SET last_message_at=NOW(),updated_at=NOW() WHERE id=?", m.ConversationID)

	return &m, nil
}
