do $$
begin

if not exists (select 1 from pg_type where typname = 'gender') then
	create type public.gender as enum ('male', 'female', 'other');
end if;

end$$;

create table if not exists public.person (
	id			bigserial primary key,
	name		varchar(1024)			 not null,
	surname		varchar(2048)			 not null,
	patronymic	varchar(2048)			 null,
	age			int						 null,
	gender		public.gender			 null,
	nationality varchar(4096)			 null,
	created     timestamptz              not null default now(),
    updated     timestamptz              not null default now(),
	removed		boolean					 not null default false
);
