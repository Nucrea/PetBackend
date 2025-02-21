create table if not exists action_tokens (
    id int primary key generated always as identity,
    user_id int references users(id),
    value text not null,
    target text not null,
    expires_at timestamp not null,
    created_at timestamp,
    updated_at timestamp

    constraint pk_action_tokens_id primary key(id),
    constraint check chk_action_tokens_target target in ('verify', 'restore')
);

create index if not exists idx_action_tokens_value on actions_tokens(value);

create or replace trigger trg_action_token_created
    before insert on action_tokens
    for each row
    when new is distinct from old
    execute function trg_proc_row_created();

create or replace trigger trg_action_token_updated
    before update on action_tokens
    for each row
    when new is distinct from old
    execute function trg_proc_row_updated();