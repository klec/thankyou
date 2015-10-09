<?php
/**
 * Created by PhpStorm.
 * User: klec
 * Date: 10/7/15
 * Time: 2:29 PM
 */

function prepareReviews($db){
    $html="";
    $page = array();
    $reviews = $db->query("select * from `reviews` LEFT JOIN `persons` on `persons`.`id`=`reviews`.`slave_id`", MYSQLI_ASSOC);
    foreach($reviews->fetch_all() as $key=>$review){
        $renderedReview = renderReview($review);
        if($key>2){

			$html.="<li>".$renderedReview."</li>\n";
		} else {
            $page["Best".($key+1)] = $renderedReview;
        }
    }
    $page["others"] = $html;
	return $page;
}

function renderReview($review){
    return $review[6]." ".$review[3];
}

function slavesOptions($db){
    $list="";
    $persons = $db->query("select * from `persons`", MYSQLI_ASSOC);
    foreach($persons->fetch_all() as $person){
        $list .= "<option value=\"$person[0]\" >$person[1]</option>";
    };

    return $list;
}

$page = prepareReviews($db);
