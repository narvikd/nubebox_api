create table testfiles (
    id uuid default gen_random_uuid() not null primary key,
    filename varchar(100) not null,
    contents bytea        not null,
    chunk_num bigint       not null
);
create unique index testfiles_filename_chunk_num_uindex on testfiles (filename, chunk_num);