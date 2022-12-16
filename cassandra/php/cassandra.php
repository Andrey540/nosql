<?php

require 'php-cassandra/php-cassandra.php';

$nodes = [
	[
		'host'		=> '172.22.0.2',
		'port'		=> 9042
	],
	[
		'host'		=> '172.22.0.3',
		'port'		=> 9042
	],
	[
		'host'		=> '172.22.0.4',
		'port'		=> 9042,
	],
];

// Create a connection.
$connection = new Cassandra\Connection($nodes, 'test_keyspace');

//Connect
try
{
	$connection->connect();
}
catch (Cassandra\Exception $e)
{
	echo 'Caught exception: ',  $e->getMessage(), "\n";
	exit;//if connect failed it may be good idea not to continue
}

// Set consistency level for farther requests (default is CONSISTENCY_ONE)
$connection->setConsistency(0x0004);

$time = microtime(true);
$memory = memory_get_usage();

$statement = $connection->queryAsync('SELECT * FROM "user"', []);
$response = $statement->getResponse();
$result = $response->fetchAll();

var_dump('time: ' . (microtime(true) - $time) . ' seconds');
var_dump('memory: ' . ((memory_get_usage() - $memory) / 1024 / 1024) . ' Mb');