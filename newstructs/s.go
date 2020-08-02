package newstructs

//this is the exported struct
type FMatrix struct {
	Row    int
	Tables [][]float64
}
type Ops interface {
	Multiply(FMatrix) FMatrix
}

// new data type
type Queue struct {
	M []FMatrix
}

//Multiply
func (matrix FMatrix) Multiply(other FMatrix) FMatrix {
	var result_matrix FMatrix
	result_matrix.Row = matrix.Row

	a := make([][]float64, result_matrix.Row)
	for i := range a {
		a[i] = make([]float64, result_matrix.Row)
	}

	for i := 0; i < matrix.Row; i++ {
		for j := 0; j < matrix.Row; j++ {
			for k := 0; k < matrix.Row; k++ {
				a[i][j] += matrix.Tables[i][k] * other.Tables[k][j]
			}
		}
	}
	result_matrix.Tables = a
	return result_matrix
}
