package main

import(
	"flag"
	"io"
	"log"
)
type InterceptorFunc func(args []string)
type Command struct{
	name string
	UsageLine string
	Short string
	Long string
	Flag flag.FlagSet
	result *io.Writer
	Interceptor InterceptorFunc 
	IsRunning bool
	Run func(cmd *Command, args []string) int
}

// Name: return command long name
func (c *Command) Name(){
	log.Println(c.UsageLine) 
}

func (c *Command) Usage(){
	log.Println(c.UsageLine) 
}

func (c *Command) Options() map[string]string {
	options := make(map[string]string)
	c.Flag.VisitAll(func(f *flag.Flag) {
		defaultVal := f.DefValue
		if len(defaultVal) > 0 {
			options[f.Name+"="+defaultVal] = f.Usage
		} else {
			options[f.Name] = f.Usage
		}
	})
	return options
}