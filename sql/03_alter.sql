use isucari;

UPDATE items i INNER JOIN categories c c.id = i.category_id SET i.parent_category_id = c.parent_id;
