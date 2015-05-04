package stats

import (
	"math"
)

func GammaPDF(alpha, beta float64) func(float64) float64 {
	return func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return (math.Pow(beta, alpha) / math.Gamma(alpha)) * math.Pow(x, alpha-1) * math.Exp(-beta*x)
	}
}

func GammaPDFAt(alpha, beta, x float64) float64 {
	pdf := GammaPDF(alpha, beta)
	return pdf(x)
}
