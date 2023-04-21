create function next_prev_id(current_id integer)
    returns TABLE(next_id integer, prev_id integer)
    language plpgsql
as
$$
begin
    select id into next_id from projects where id > current_id order by id limit 1;
    if not found then
        select id into next_id from projects order by id limit 1;
    end if;
    select id into prev_id from projects where id < current_id order by id desc limit 1;
    if not found then
        select id into prev_id from projects order by id desc limit 1;
    end if;
    return next;
end;
$$;