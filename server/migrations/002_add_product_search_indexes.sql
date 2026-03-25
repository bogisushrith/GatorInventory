-- Add indexes to improve product search and filtering performance
CREATE INDEX IF NOT EXISTS idx_products_name_lower ON products (LOWER(name));
CREATE INDEX IF NOT EXISTS idx_products_category_lookup ON products(category);
