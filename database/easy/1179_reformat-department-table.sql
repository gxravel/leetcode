-- 1179. Reformat Department Table

-- Table: Department

-- +---------------+---------+
-- | Column Name   | Type    |
-- +---------------+---------+
-- | id            | int     |
-- | revenue       | int     |
-- | month         | varchar |
-- +---------------+---------+
-- (id, month) is the primary key of this table.
-- The table has information about the revenue of each department per month.
-- The month has values in ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"].

-- Write an SQL query to reformat the table such that there is a department id column and a revenue column for each month.

-- The query result format is in the following example:

-- Department table:
-- +------+---------+-------+
-- | id   | revenue | month |
-- +------+---------+-------+
-- | 1    | 8000    | Jan   |
-- | 2    | 9000    | Jan   |
-- | 3    | 10000   | Feb   |
-- | 1    | 7000    | Feb   |
-- | 1    | 6000    | Mar   |
-- +------+---------+-------+

-- Result table:
-- +------+-------------+-------------+-------------+-----+-------------+
-- | id   | Jan_Revenue | Feb_Revenue | Mar_Revenue | ... | Dec_Revenue |
-- +------+-------------+-------------+-------------+-----+-------------+
-- | 1    | 8000        | 7000        | 6000        | ... | null        |
-- | 2    | 9000        | null        | null        | ... | null        |
-- | 3    | null        | 10000       | null        | ... | null        |
-- +------+-------------+-------------+-------------+-----+-------------+

-- Note that the result table has 13 columns (1 for the department id + 12 for the months).

SELECT
  id,
  max(case month when "Jan" then revenue else null end) as Jan_Revenue,
  max(case month when "Feb" then revenue else null end) as Feb_Revenue,
  max(case month when "Mar" then revenue else null end) as Mar_Revenue,
  max(case month when "Apr" then revenue else null end) as Apr_Revenue,
  max(case month when "May" then revenue else null end) as May_Revenue,
  max(case month when "Jun" then revenue else null end) as Jun_Revenue,
  max(case month when "Jul" then revenue else null end) as Jul_Revenue,
  max(case month when "Aug" then revenue else null end) as Aug_Revenue,
  max(case month when "Sep" then revenue else null end) as Sep_Revenue,
  max(case month when "Oct" then revenue else null end) as Oct_Revenue,
  max(case month when "Nov" then revenue else null end) as Nov_Revenue,
  max(case month when "Dec" then revenue else null end) as Dec_Revenue
  FROM
  Department
  group by id;