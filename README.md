### Overview

The solution implements a simple HTTPS server to fetch the promotion details by its ID.

The server is written in Go. `chi` is used as a router and `sqlx` library is used to scan rows into a `struct`.

SSL-certificate is self-signed. Thus, `curl` should be executed with `-k` parameter to ignore SSL errors.

`MySQL` is used as a storage. 

To store objects from the `.csv` file the `LOAD DATA INFILE` command is used. In the current example this command loads records from the local file. Otherwise, the file may be copied to the machine, where `MySQL`-server is running, and then the command could be executed. Also this command transforms `id` to the upper-case string and `expiration_date` to UTC. 

### Usage

To run the example:

```bash
docker-compose up -d
```

To enable loading of a local file please connect to `MySQL`:

```bash
docker run --network="host" -v /folder/with/csv-file:/data -it --rm mysql mysql -h 127.0.0.1 -P 3306 --protocol=tcp --local-infile=1 -uroot -p
```

The root password is `test`.

And then:

```sql
mysql> SET GLOBAL local_infile=1;
mysql> quit
```

Connect to `MySQL` again:

```bash
docker run --network="host" -v /folder/with/csv-file:/data -it --rm mysql mysql -h 127.0.0.1 -P 3306 --protocol=tcp --local-infile=1 -uroot -p
```

To load data from file the following commands are used:

```sql
mysql> use promotions;
mysql> LOAD DATA LOCAL INFILE '/data/promotions.csv' INTO TABLE promotions FIELDS TERMINATED BY "," LINES TERMINATED BY "\n" IGNORE 1 LINES (@in_id,price,@exp_date) SET id=UPPER(@in_id), expiration_date=CONVERT_TZ(SUBSTR(@exp_date, 1, 19), INSERT(SUBSTR(@exp_date, 21, 5), 4, 1, ':'), '+00:00');
```

Note that the filename in this example is `promotions.csv`. 

To get the promotion details:

```bash
curl -k https://localhost:1321/promotions/EC9D76A0-9A3F-4A59-9C1A-7C388775A896
{"id":"EC9D76A0-9A3F-4A59-9C1A-7C388775A896","price":76.044174,"expiration_date":"2018-06-08 04:12:54"}
```

### Further enhancements

The type of `id` column could be changed from `CHAR(36)` to `BINARY(16)` to reduce space usage. In this case the `SELECT` statement in our service should be changed to:

```sql
SELECT BIN_TO_UUID(),price,expiration_date FROM promotions WHERE id=UUID_TO_BIN(?)
```

Also we can try to use the `nginx` web-server with `handlersocket_json_module` (https://github.com/nginx-modules/ngx_http_handlersocket_json_module). In this case we don't even need service to fetch the data from the database.

Another approach is to use separate service/tool to parse `.csv` file and write data to the storage. Multiple goroutines can be used for writing. To speed up the process multi-row inserts can be used with some tips: https://support.tigertech.net/mysql-large-inserts . 

If we have a sufficient amount of RAM, the in-memory key-value store like `Redis` can be considered as a storage.

### Additional considerations

- Please consider the `.csv` file to be very big (billions of entries) - how would your application perform?
- How would your application perform in peak periods (millions of request per minute)?
- Every new file is immutable, that is, you should erase and write the whole storage.

