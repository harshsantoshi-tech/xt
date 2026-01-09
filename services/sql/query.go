package sql

const GET_USERS_FROM_EMAIL = `SELECT id, username, email, password_hash FROM users WHERE email = ? LIMIT 1`

const UPDATE_PASSWORD_BY_EMAIL = `UPDATE users SET password_hash = ? WHERE email = ?`

const INSERT_USERS = "INSERT INTO users (username , email, password_hash, created_at) VALUES (?,?, ?, ?)"

const GET_CHAT_LIST = `
    SELECT 
        r.id AS room_id,
        r.name AS room_name,
        r.is_group,
        r.last_message_at,
        m.content AS last_message,
        m.message_type,
        m.sender_id AS last_sender_id,
        u.username AS last_sender_name,
        (SELECT COUNT(*) FROM messages 
         WHERE room_id = r.id 
         AND id > rm.last_read_message_id 
         AND sender_id != ?) AS unread_count
    FROM rooms r
    JOIN room_members rm ON r.id = rm.room_id
    LEFT JOIN messages m ON m.room_id = r.id AND m.created_at = r.last_message_at
    LEFT JOIN users u ON m.sender_id = u.id
    WHERE rm.user_id = ?
    ORDER BY r.last_message_at DESC;
`


const GET_CHAT_LIST_PAGINATED = `
    SELECT 
        r.id AS room_id,
        r.name AS room_name,
        r.is_group,
        r.last_message_at,
        m.content AS last_message,
        m.message_type,
        m.sender_id AS last_sender_id,
        (SELECT username FROM users WHERE id = m.sender_id) AS last_sender_name,
        (SELECT COUNT(*) FROM messages 
         WHERE room_id = r.id 
         AND id > rm.last_read_message_id 
         AND sender_id != ?) AS unread_count,
        u_other.id AS other_user_id,
        u_other.username AS other_user_name
    u_other.last_seen_at AS other_last_seen
    FROM rooms r
    JOIN room_members rm ON r.id = rm.room_id
    LEFT JOIN room_members rm_other ON rm_other.room_id = r.id AND rm_other.user_id != rm.user_id AND r.is_group = 0
    LEFT JOIN users u_other ON rm_other.user_id = u_other.id
    LEFT JOIN messages m ON m.room_id = r.id AND m.created_at = r.last_message_at
    WHERE rm.user_id = ?
    ORDER BY r.last_message_at DESC
    LIMIT ? OFFSET ?;
`

const GET_MESSAGES_BY_ROOM = `
    SELECT 
        id, 
        room_id, 
        sender_id, 
        content, 
        message_type, 
        created_at 
    FROM messages 
    WHERE room_id = ? AND is_deleted = 0
    ORDER BY created_at DESC 
    LIMIT ? OFFSET ?;
`