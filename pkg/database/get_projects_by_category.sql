create function get_projects_by_category(pcategory text)
    returns TABLE(id integer, name text, thumbnail text)
    language sql
as
$$
select p.id, p.name, coalesce((select url from pictures where project_id = p.id order by is_thumbnail desc nulls last limit 1)) from projects p where p.category = pcategory order by p.id desc;
$$;

alter function get_projects_by_category(text) owner to doadmin;

