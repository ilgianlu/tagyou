package conf

// AI url:
// http endpoint where an openAI api compatible model is responding
var AI_URL string = "http://localhost:11434/v1"

// AI model to use
var AI_MODEL string = "gemma3:4b-it-qat"

// should be a useless key for AI_URL
var AI_API_KEY string = "useless"

// folder where to save AI blueprint for device traffic
var AI_BLUEPRINT_PATH string = "./blueprint"
