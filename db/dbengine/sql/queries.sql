-- name: GetFileByID :one
select contents
from testfiles
where filename = $1 and chunk_num = $2
limit 1;

-- name: GetFileChunkCount :one
select count(chunk_num)
from testfiles
where filename = $1
limit 1;

-- name: InsertFile :execresult
INSERT INTO testfiles (id, filename, contents, chunk_num)
VALUES (DEFAULT, $1, $2, $3);

-- name: GetAllFileNames :many
select distinct filename from testfiles;

-- name: DeleteFile :execresult
delete
from testfiles
where filename = $1;

-- name: FileExists :one
select exists(
    select 1
    from testfiles where filename = $1
);
