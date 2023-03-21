create function get_project(pid integer)
    returns TABLE(id integer, name text, description text, video text, type text, category text, mentor text, has_thumbnail boolean, demo text, github text)
    language sql
as
$$
select p.id, p.name, p.description, p.video_link as video, p.type, p.category, p.mentor, p.has_thumbnail, p.demo_link as demo, p.github_link as github from projects p where p.id = pid;
$$;

alter function get_project(integer) owner to doadmin;

