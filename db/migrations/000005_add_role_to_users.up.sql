ALTER TABLE "users" ADD COLUMN "role" varchar NOT NULL DEFAULT 'depositor';


INSERT INTO "users"(
	username, hashed_password, full_name, email, password_changed_at, created_at, is_email_verified, role)
	VALUES ('my_user', '$2a$10$vLRcCeRqMsr7QA71cQ3ChO9lHtTdwJK/iVokChIpafDlEkyy3foS6', 'Full Name', 'full.name@email.com', now(), now(), true, 'none_yet');