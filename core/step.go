package core

import (
	"fmt"

	"github.com/Ether-Security/leviathan/libs"
	"github.com/Ether-Security/leviathan/utils"
	"github.com/antonmedv/expr"
)

func (r *Runner) RunStep(step *libs.Step) {
	// If there is a custom source
	if step.Source != "" {
		lines := utils.ReadFile(r.ParseString(step.Source))
		if len(lines) == 0 {
			utils.Logger.Error().Msg("Source not found or empty")
			return
		}
		for id, line := range lines {
			var stepClone libs.Step
			stepClone = *step
			stepClone.Source = ""

			r.Params["source"] = line
			r.Params["source_safe"] = utils.CleanPath(line)
			r.Params["source_id"] = fmt.Sprintf("%d", id)

			r.RunStep(&stepClone)
		}
	}

	// Check requirements, skip step if requirements not met
	if !r.CheckRequirements(step.Requirements) {
		return
	}

	// If conditions are validated
	if r.CheckConditions(step.Conditions) {
		// Execute commands
		r.runCommands(step.Commands)

		// Execute scripts
		r.runScripts(step.Scripts)
	} else {
		// Execute commands
		r.runCommands(step.RCommands)

		// Execute scripts
		r.runScripts(step.RScripts)
	}

}

func (r *Runner) CheckRequirements(requirements []string) bool {
	for _, requirement := range requirements {
		parsedReq := r.ParseString(requirement)

		if utils.FileExists(parsedReq) {
			utils.Logger.Debug().Str("requirement", parsedReq).Msg("File found")
			continue
		}
		if utils.FolderExists(parsedReq) {
			utils.Logger.Debug().Str("requirement", parsedReq).Msg("Folder found")
			continue
		}
		if utils.BinaryExists(parsedReq) {
			utils.Logger.Debug().Str("requirement", parsedReq).Msg("Binary found")
			continue
		}
		utils.Logger.Fatal().Str("requirement", parsedReq).Msg("Requirement not found")
		return false
	}
	return true
}

func (r *Runner) CheckConditions(conditions []string) bool {
	for _, condition := range conditions {
		parsedCondition := r.ParseString(condition)
		utils.Logger.Debug().Str("script", parsedCondition).Msg("Check condition")
		ret, err := expr.Eval(parsedCondition, r.DSL)
		if err != nil {
			utils.Logger.Warn().Str("condition", parsedCondition).Msg(err.Error())
			return false
		}
		if ret.(bool) == false {
			utils.Logger.Debug().Msg("Condition not met")
			return false
		}
	}
	return true
}

func (r *Runner) runCommands(commands []string) {
	for _, cmd := range commands {
		parsedCmd := r.ParseString(cmd)
		utils.Logger.Debug().Str("command", parsedCmd).Msg("Execute command")
		_, err := utils.ExecCmd(parsedCmd, r.Workspace, r.Options.Environment.Binaries, r.Options.Log.Debug)
		if err != nil {
			utils.Logger.Error().Msg(err.Error())
		}
	}
}

func (r *Runner) runScripts(scripts []string) {
	for _, cmd := range scripts {
		parsedCmd := r.ParseString(cmd)
		utils.Logger.Debug().Str("script", parsedCmd).Msg("Execute script")
		_, err := expr.Eval(parsedCmd, r.DSL)
		if err != nil {
			utils.Logger.Warn().Str("script", parsedCmd).Msg(err.Error())
			continue
		}
	}
}
