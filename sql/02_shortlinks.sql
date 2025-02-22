create table if not exists shortlinks (
    id int generated always as identity,
    url text not null,
    expires_at timestamp not null,
    created_at timestamp,
    updated_at timestamp
);

create or replace trigger trg_shortlink_created
    before insert on shortlinks
    for each row
    execute function trg_proc_row_created();

create or replace trigger trg_shortlink_updated
    before update on shortlinks
    for each row 
    when (new is distinct from old)
    execute function trg_proc_row_updated();