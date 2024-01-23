do $$
begin

if exists (select 1 from pg_type where typname = 'gender') then
	drop type public.gender;
end if;

end$$;

drop table if exists public.person;

