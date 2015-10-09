<?php
/**
 * Created by PhpStorm.
 * User: klec
 * Date: 10/7/15
 * Time: 2:29 PM
 */

$db = new mysqli('', 'root', 'ryurik', 'thankyou', 0,  '/Applications/MAMP/tmp/mysql/mysql.sock');
if (!$db) {
    die("Error: Unable to connect to MySQL." . PHP_EOL);
}