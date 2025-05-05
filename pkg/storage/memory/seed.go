package memory

import (
	"meche/pkg/models"
	"meche/pkg/storage"
)

// SeedData initializes the storage with sample data
func SeedData(orgStorage storage.OrganizationStorage, projectStorage storage.ProjectStorage) {
	// Create a sample organization
	org := models.NewOrganization("Acme Corp", "A leading technology company focused on innovation and excellence")
	orgStorage.CreateOrganization(org)

	// Create three sample projects for the organization
	projects := []*models.Project{
		models.NewProject(
			"Website Redesign",
			"Modernize the company website with a new design and improved user experience",
			org.ID,
		),
		models.NewProject(
			"Mobile App Development",
			"Build a cross-platform mobile application for customer engagement",
			org.ID,
		),
		models.NewProject(
			"Data Analytics Platform",
			"Develop a platform for analyzing customer data and generating insights",
			org.ID,
		),
	}

	// Add the projects to storage
	for _, project := range projects {
		projectStorage.CreateProject(project)
	}
}
