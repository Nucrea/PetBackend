create or replace function trg_proc_row_updated()
returns trigger as $$
begin
    if new is distinct from old then
        new.updated_at = now();
    end if;
    return new; 
end;
$$ language plpgsql;

create or replace function trg_proc_row_created()
returns trigger as $$
begin
    new.created_at = now();
	new.updated_at = now();
    return new;
end;
$$ language plpgsql;