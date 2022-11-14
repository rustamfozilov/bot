create table my_users(
    id bigserial primary key,
    user_id bigint unique,
    mail_login text,
    mail_password  text,
    mail_service text,
    total_msg_count bigint,
    unseen_msg_count bigint,
    created_at timestamptz default current_timestamp
);

drop table my_users;