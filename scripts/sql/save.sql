SELECT p.UUID,p.Name,p.Fnam,p.Mnam,p.Snam,p.Sex,p.DOB,p.DOD,s.bitchId
FROM people p
INNER JOIN partners s ON p.UUID = s.butchId
ORDER BY p.UUID
