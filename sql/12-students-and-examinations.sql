-- https://leetcode.com/problems/students-and-examinations/description/

SELECT student_subject.student_id,
       student_subject.student_name,
       student_subject.subject_name,
       COUNT(e.subject_name) AS attended_exams
FROM (
    SELECT student_id, student_name, subject_name
    FROM (
        SELECT s.student_id, s.student_name, sub.subject_name
        FROM Students s
        join Subjects sub
        on 1=1
    ) AS pairs
) AS student_subject
LEFT JOIN Examinations e
  ON student_subject.student_id = e.student_id
 AND student_subject.subject_name = e.subject_name
GROUP BY student_subject.student_id, student_subject.student_name, student_subject.subject_name
ORDER BY student_subject.student_id, student_subject.subject_name;


 