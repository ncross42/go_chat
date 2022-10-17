CREATE OR REPLACE FUNCTION new_message_event_notif() RETURNS TRIGGER AS $$
    DECLARE
        sender_username VARCHAR;
        notif json;
    BEGIN
        SELECT username INTO sender_username
            FROM users WHERE id = NEW.sender_id;

        notif = json_build_object(
            'sender', sender_username,
            'senderID', NEW.sender_id,
            'text', NEW.text,
            'sentOn', NEW.sent_on
        );
        PERFORM pg_notify('new_message_' || NEW.chatroom_id, notif::text);
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER new_message_event_notif_trigger
    AFTER INSERT ON messages
    FOR EACH ROW EXECUTE PROCEDURE new_message_event_notif();