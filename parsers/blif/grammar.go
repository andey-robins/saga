package blif

type BlifFile struct {
	Comment string  `@Comment? NewLine?`
	Header  *Header `@@`
	Gates   []*Gate `@@* ".end" NewLine?`
}

type Header struct {
	ModelName  string  `".model" @Ident NewLine`
	InputList  *IOList `".inputs" @@ NewLine`
	OutputList *IOList `".outputs" @@ NewLine`
}

type IOList struct {
	Nodes []string `@Ident*`
}

type Gate struct {
	NotGate *NotGate `  ".gate" "NOT" @@ NewLine`
	NorGate *NorGate `| ".gate" "NOR" @@ NewLine`
}

type NotGate struct {
	InputName  string `"A" "=" @Ident`
	OutputName string `"Y" "=" @Ident`
}

type NorGate struct {
	InputOneName string `"A" "=" @Ident`
	InputTwoName string `"B" "=" @Ident`
	OutputName   string `"Y" "=" @Ident`
}
