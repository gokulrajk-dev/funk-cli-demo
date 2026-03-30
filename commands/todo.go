package commands

import (
	"context"
	"fmt"
	"funk/sqldb"
	"github.com/urfave/cli/v3"
)


func Todos() *cli.Command {
	return &cli.Command{
		Name:   "todo",
		Usage:  "using todo to store the task via terminal",
		Action: Task,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "add",
				Usage: "use to add new task",
				Aliases: []string{"a","i"},
			},
			&cli.IntFlag{
				Name:  "del",
				Usage: "use to delete a task",
				Aliases: []string{"d","rm"},
			},
			&cli.BoolFlag{
				Name:  "task",
				Usage: "use to show all task",
				Aliases: []string{"t"},
			},
			&cli.IntFlag{
				Name:  "done",
				Usage: "use to update task status",
				Aliases: []string{"c"},
			},
		},
	}
}

func Task(ctx context.Context, cmd *cli.Command) error {

	sqldb.InitTable()

	switch {
	case cmd.IsSet("add"):
		task := cmd.String("add")
		sqldb.AddTask(task)
		
	case cmd.Bool("task"):
		sqldb.ListTasks()

	case cmd.IsSet("del"):
		index := cmd.Int("del")
		sqldb.Delete_task(index)

	case cmd.IsSet("done"):
		index:=cmd.Int("done")
		fmt.Println("did you complete the task (y/n)")
		var complete string
		fmt.Scan(&complete)
		if len(complete) > 0 && complete[0] == 'y'{
		sqldb.Update_Task(index)
		}else{
			fmt.Println("complete the task and afterwards update")
		}
		

	default:
		return fmt.Errorf("please specify one of: --add, --del, or --task")
	}

	return nil

}