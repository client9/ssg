package shortcode

// Entity is a HandlerFunc that converts its single argument to an HTML entity.
//
//	$ent[across]  -->  &across;
//	$ent[#1245]   -->  &#1245;
func Entity(_ *Context, args []string, _ string) string {
	if len(args) != 1 {
		return ""
	}
	return "&" + args[0] + ";"
}
