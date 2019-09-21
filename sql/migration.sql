DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `seller_id` bigint NOT NULL,
  `buyer_id` bigint NOT NULL,
  `trans_status` enum('wait_shipping', 'wait_done', 'done') NOT NULL,
  `item_id` bigint NOT NULL UNIQUE,
  `item_name` varchar(191) NOT NULL,
  `item_price` int unsigned NOT NULL,
  `item_description` text NOT NULL,
  `item_category_id` int unsigned NOT NULL,
  `item_root_category_id` int unsigned NOT NULL,

  `ship_status` enum('initial', 'wait_pickup', 'shipping', 'done') NOT NULL,
  `reserve_id` varchar(191) NOT NULL,
  `reserve_time` bigint NOT NULL,
  `to_address` varchar(191) NOT NULL,
  `to_name` varchar(191) NOT NULL,
  `from_address` varchar(191) NOT NULL,
  `from_name` varchar(191) NOT NULL,
  `img_binary` mediumblob NOT NULL,

  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `trans_updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

INSERT INTO `transactions`(id, seller_id, buyer_id, trans_status, item_id, item_name, item_price, item_description, item_category_id, item_root_category_id, ship_status, reserve_id, reserve_time, to_address, to_name, from_address, from_name, img_binary, created_at, updated_at, trans_updated_at)
SELECT t.id, t.seller_id, t.buyer_id, t.status, t.item_id, t.item_name, t.item_price, t.item_description, t.item_category_id, t.item_root_category_id, s.status, s.reserve_id, s.reserve_time, s.to_address, s.to_name, s.from_address, s.from_name, s.img_binary, s.created_at, s.updated_at, t.updated_at FROM `transaction_evidences` t INNER JOIN `shippings` s ON t.`id` = s.transaction_evidence_id;

DROP TABLE IF EXISTS `transaction_evidences`;
DROP TABLE IF EXISTS `shippings`;
