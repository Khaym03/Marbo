package expansion

import "github.com/Khaym03/Marbo/internal/domain"

type IntentSeedSet struct {
	IntentID domain.IntentID
	Concepts []string
}

func GetSeeds() []IntentSeedSet {
	return []IntentSeedSet{
		{
			IntentID: "INTENT_REQUISITOS_INGRESO",
			Concepts: []string{"inscripción", "ingresar", "estudiar", "documentos", "requisitos", "nuevo ingreso"},
		},
		{
			IntentID: "INTENT_FECHAS_REGULARES",
			Concepts: []string{"inscripción", "semestre", "fechas", "materias", "reinscripción"},
		},
		{
			IntentID: "INTENT_SAIA_ACCESO",
			Concepts: []string{"SAIA", "plataforma", "acceso", "usuario", "contraseña"},
		},
		{
			IntentID: "INTENT_PROBLEMA_ACTA",
			Concepts: []string{"lista", "profesor", "acta", "inscripción", "materia"},
		},
		{
			IntentID: "INTENT_INFO_SEDES",
			Concepts: []string{"sede", "dirección", "ubicación", "edificio", "campus"},
		},
		{
			IntentID: "INTENT_ACTIVIDADES_EXTRACURRICULARES",
			Concepts: []string{"deportes", "cultura", "actividades", "estudiantes", "eventos"},
		},
	}
}
