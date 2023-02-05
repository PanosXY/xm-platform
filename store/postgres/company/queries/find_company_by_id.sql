SELECT id, name, description, employees_amount, registered, type
FROM companies
WHERE id = $1
