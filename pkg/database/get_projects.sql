create function get_projects()
    returns TABLE(id integer, name text, thumbnail text, category text)
    language sql
as
$$
select p.id, p.name, coalesce((select url from pictures where project_id = p.id order by is_thumbnail desc nulls last limit 2)), p.category from projects p order by p.id;
$$;

alter function get_projects() owner to doadmin;

