const asyncHandler = require("express-async-handler");
const { processDBRequest } = require("../../utils");
const {
  getAllStudents,
  addNewStudent,
  getStudentDetail,
  setStudentStatus,
  updateStudent,
} = require("./students-service");

const handleGetAllStudents = asyncHandler(async (req, res) => {
  //write your code
});

const handleAddStudent = asyncHandler(async (req, res) => {
  //write your code
});

const handleUpdateStudent = asyncHandler(async (req, res) => {
  //write your code
});

const handleGetStudentDetail = asyncHandler(async (req, res) => {
  const { id } = req.params;
  const query = `SELECT *
                   FROM students
                  WHERE id = $1;`;

  const queryParams = [id];

  const { rows } = await processDBRequest({ query, queryParams });

  if (rows && rows.length == 0) {
    res.status(404).json({ error: "no student with this id" });
    return;
  }
  res.json(rows[0]);
});

const handleStudentStatus = asyncHandler(async (req, res) => {
  //write your code
});

module.exports = {
  handleGetAllStudents,
  handleGetStudentDetail,
  handleAddStudent,
  handleStudentStatus,
  handleUpdateStudent,
};
