package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"go/format"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	mdt "github.com/Gurvan/melee-data-tools"
	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/descriptor"
	"github.com/Gurvan/melee-data-tools/fighter"
	"github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/model"
	"github.com/davecgh/go-spew/spew"
)

type section struct {
	key   string
	title string
	body  func() string
}

type reviewAnimation struct {
	Index  int
	Offset int64
	Size   int
	Root   string
	Data   mdt.AnimationData
}

type actionAnimationRef struct {
	Index int
	Name  string
	Size  uint32
}

func main() {
	sectionKey := flag.String("section", "", "print one section and exit")
	listSections := flag.Bool("list", false, "list available sections and exit")
	animationPath := flag.String("animations", "", "parse a fighter animation bundle and list all animations, such as PlMsAJ.dat")
	writeSpecial := flag.Bool("write-special", false, "create or overwrite fighter/attributes/<id>.go and update special.go")
	specialID := flag.String("special-id", "", "Go type/file id for -write-special, such as Fc, Pr, or Ca")
	specialType := flag.String("special-type", "auto", "field type for generated unknown attributes: auto, uint32, int32, float32, or Hex32")
	force := flag.Bool("force", false, "overwrite generated special attribute files without prompting")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage: %s [options] <fighter.dat>

Review parsed fighter data without editing a scratch test.

Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	path := expandHome(flag.Arg(0))
	var fighterFile mdt.FighterFile
	fighterData, desc, err := fighterFile.ReadFromFile(path)
	if err != nil {
		log.Fatalf("Could not parse %s as FighterFile: %v", path, err)
	}

	if *writeSpecial {
		if err := writeSpecialAttributes(path, fighterData, desc, *specialID, *specialType, *force); err != nil {
			log.Fatal(err)
		}
		return
	}

	var animations []reviewAnimation
	if path, required := resolveAnimationPath(path, *animationPath); path != "" {
		animations, err = readAnimationBundle(path)
		if err != nil {
			if required {
				log.Fatalf("Could not parse %s as animation bundle: %v", path, err)
			}
			log.Printf("Could not parse inferred animation bundle %s: %v", path, err)
		}
	}

	if *sectionKey == "" && !*listSections {
		if err := maybeOfferWriteSpecial(path, fighterData, desc, *specialType); err != nil {
			log.Fatal(err)
		}
	}

	sections := buildSections(path, fighterData, desc, animations)
	if *listSections {
		for _, section := range sections {
			fmt.Printf("%-16s %s\n", section.key, section.title)
		}
		return
	}

	if *sectionKey != "" {
		section, ok := findSection(sections, *sectionKey)
		if !ok {
			log.Fatalf("Unknown section %q. Run with -list to see options.", *sectionKey)
		}
		printSection(section)
		return
	}

	runMenu(sections)
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func resolveAnimationPath(fighterPath string, animationPath string) (string, bool) {
	if animationPath != "" {
		return expandHome(animationPath), true
	}

	inferred := inferAnimationPathFromFighterPath(fighterPath)
	if inferred == "" || !fileExists(inferred) {
		return "", false
	}
	return inferred, false
}

func inferAnimationPathFromFighterPath(fighterPath string) string {
	dir := filepath.Dir(fighterPath)
	name := filepath.Base(fighterPath)
	matches := regexp.MustCompile(`(?i)^(Pl[A-Za-z0-9]{2})\.dat$`).FindStringSubmatch(name)
	if len(matches) != 2 {
		return ""
	}
	return filepath.Join(dir, matches[1]+"AJ.dat")
}

func buildSections(path string, data mdt.FighterData, desc descriptor.Descriptor, animations []reviewAnimation) []section {
	return []section{
		{
			key:   "summary",
			title: "Summary",
			body: func() string {
				return renderSummary(path, data, desc)
			},
		},
		{
			key:   "descriptor",
			title: "Descriptor",
			body: func() string {
				return dump(desc)
			},
		},
		{
			key:   "header",
			title: "Descriptor Header",
			body: func() string {
				return dump(desc.Header)
			},
		},
		{
			key:   "roots",
			title: "Roots",
			body: func() string {
				return dump(desc.Roots)
			},
		},
		{
			key:   "relocation",
			title: "Relocation Table",
			body: func() string {
				return renderRelocation(desc)
			},
		},
		{
			key:   "common",
			title: "Common Attributes",
			body: func() string {
				return dumpPtr(data.AttributesCommon.ValuePtr)
			},
		},
		{
			key:   "special",
			title: "Special Attributes",
			body: func() string {
				return renderSpecialAttributes(data, desc)
			},
		},
		{
			key:   "model-params",
			title: "Model Params",
			body: func() string {
				return dumpPtr(data.ModelParams.ValuePtr)
			},
		},
		{
			key:   "shield-pose",
			title: "Shield Pose",
			body: func() string {
				return renderShieldPose(data)
			},
		},
		{
			key:   "dynamics",
			title: "Dynamic Bones",
			body: func() string {
				return dumpPtr(data.Dynamics.ValuePtr)
			},
		},
		{
			key:   "hurtboxes",
			title: "Hurtboxes",
			body: func() string {
				return dumpPtr(data.Hurtboxes.ValuePtr)
			},
		},
		{
			key:   "ecb",
			title: "Environment Collision Box",
			body: func() string {
				return dumpPtr(data.ECB.ValuePtr)
			},
		},
		{
			key:   "jostle",
			title: "Jostle Box",
			body: func() string {
				return dumpPtr(data.JostleBox.ValuePtr)
			},
		},
		{
			key:   "items",
			title: "Items",
			body: func() string {
				return renderItems(data)
			},
		},
		{
			key:   "actions",
			title: "Actions",
			body: func() string {
				return renderActions(data)
			},
		},
		{
			key:   "animations",
			title: "Animations",
			body: func() string {
				return renderAnimations(data, animations)
			},
		},
		{
			key:   "subactions",
			title: "Subactions",
			body: func() string {
				return renderSubactions(data)
			},
		},
		{
			key:   "model",
			title: "Model Root",
			body: func() string {
				return dumpPtr(data.Model.ValuePtr)
			},
		},
	}
}

func renderSummary(path string, data mdt.FighterData, desc descriptor.Descriptor) string {
	var b strings.Builder
	firstRoot, err := desc.FirstRootName()
	if err != nil {
		firstRoot = "<none>"
	}
	fmt.Fprintf(&b, "File: %s\n", path)
	fmt.Fprintf(&b, "First root: %s\n", firstRoot)
	fmt.Fprintf(&b, "File size: %d\n", desc.FileSize)
	fmt.Fprintf(&b, "Relocations: %d\n", desc.RelocationCount)
	fmt.Fprintf(&b, "Roots: %d\n", desc.RootCount)
	fmt.Fprintf(&b, "Refs: %d\n", desc.RefCount)
	fmt.Fprintf(&b, "Version: %s\n", desc.Version.String())
	fmt.Fprintf(&b, "Common attributes: %s\n", yesNo(data.AttributesCommon.ValuePtr != nil))
	specialSize := blockSize(desc, data.AttributesSpecial.Offset)
	fmt.Fprintf(&b, "Special attributes: %s", yesNo(data.AttributesSpecial.ValuePtr != nil))
	if specialSize > 0 {
		fmt.Fprintf(&b, " (%d bytes, %d 4-byte values)", specialSize, specialSize/4)
	}
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Model params: %s\n", yesNo(data.ModelParams.ValuePtr != nil))
	fmt.Fprintf(&b, "Hurtboxes: %d\n", lenHurtboxes(data.Hurtboxes.ValuePtr))
	fmt.Fprintf(&b, "Items: %d\n", len(data.Items))
	fmt.Fprintf(&b, "Actions: %d\n", lenActionTable(data.ActionTable.ValuePtr))
	if data.Dynamics.ValuePtr == nil {
		fmt.Fprintln(&b, "Dynamic bones: 0 sets, 0 colliders")
	} else {
		fmt.Fprintf(&b, "Dynamic bones: %d sets, %d colliders\n", len(data.Dynamics.ValuePtr.BoneSets), len(data.Dynamics.ValuePtr.Colliders))
	}
	fmt.Fprintf(&b, "ECB: %s\n", yesNo(data.ECB.ValuePtr != nil))
	fmt.Fprintf(&b, "Jostle box: %s\n", yesNo(data.JostleBox.ValuePtr != nil))
	fmt.Fprintf(&b, "Model root: %s\n", yesNo(data.Model.ValuePtr != nil))
	return b.String()
}

func renderSpecialAttributes(data mdt.FighterData, desc descriptor.Descriptor) string {
	size := blockSize(desc, data.AttributesSpecial.Offset)
	var b strings.Builder
	fmt.Fprintf(&b, "Offset: %s\n", data.AttributesSpecial.Offset.String())
	if size > 0 {
		fmt.Fprintf(&b, "Raw block size: %d bytes\n", size)
		fmt.Fprintf(&b, "4-byte value count: %d\n", size/4)
	} else {
		fmt.Fprintf(&b, "Raw block size: unknown; this pointer is the last relocation target or was not relocated.\n")
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, dumpPtr(data.AttributesSpecial.ValuePtr))
	return b.String()
}

func writeSpecialAttributes(path string, data mdt.FighterData, desc descriptor.Descriptor, id string, fieldType string, force bool) error {
	firstRoot, err := desc.FirstRootName()
	if err != nil {
		return err
	}
	rootName := strings.TrimPrefix(firstRoot, "ftData")
	if rootName == firstRoot || rootName == "" {
		return fmt.Errorf("first root %q does not look like fighter data", firstRoot)
	}

	size := blockSize(desc, data.AttributesSpecial.Offset)
	if size == 0 {
		return fmt.Errorf("could not infer special attribute block size from relocation offset %s", data.AttributesSpecial.Offset.String())
	}
	if size%4 != 0 {
		return fmt.Errorf("special attribute block size %d is not divisible by 4", size)
	}
	if !isSupportedFieldType(fieldType) {
		return fmt.Errorf("unsupported special attribute field type %q; use auto, uint32, int32, float32, or Hex32", fieldType)
	}

	reader := bufio.NewReader(os.Stdin)
	if id == "" {
		fmt.Printf("Fighter root: %s\n", firstRoot)
		fmt.Printf("Special attribute block: %d bytes (%d 4-byte values)\n", size, size/4)
		fmt.Print("Go type/file id (for example Fc, Pr, Ca): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		id = strings.TrimSpace(input)
	}
	if err := validateSpecialID(id); err != nil {
		return err
	}

	moduleRoot, err := findModuleRoot()
	if err != nil {
		return err
	}
	attributesDir := filepath.Join(moduleRoot, "fighter", "attributes")
	attributePath := filepath.Join(attributesDir, strings.ToLower(id)+".go")
	if fileExists(attributePath) && !force {
		fmt.Printf("%s already exists. Overwrite? [y/N]: ", attributePath)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if !strings.EqualFold(strings.TrimSpace(input), "y") && !strings.EqualFold(strings.TrimSpace(input), "yes") {
			return fmt.Errorf("not overwriting %s", attributePath)
		}
	}

	fieldTypes, err := specialFieldTypes(path, data.AttributesSpecial.Offset, int(size/4), fieldType)
	if err != nil {
		return err
	}
	source, err := generateSpecialAttributeSource(id, fieldTypes)
	if err != nil {
		return err
	}
	if err := os.WriteFile(attributePath, source, 0o644); err != nil {
		return err
	}
	if err := updateSpecialSwitch(filepath.Join(attributesDir, "special.go"), rootName, id); err != nil {
		return err
	}

	fmt.Printf("Wrote %s with %d unknown attributes.\n", attributePath, size/4)
	fmt.Printf("Updated fighter/attributes/special.go: ftData%s -> %s.\n", rootName, id)
	return nil
}

func maybeOfferWriteSpecial(path string, data mdt.FighterData, desc descriptor.Descriptor, fieldType string) error {
	firstRoot, err := desc.FirstRootName()
	if err != nil {
		return err
	}
	rootName := strings.TrimPrefix(firstRoot, "ftData")
	if rootName == firstRoot || rootName == "" {
		return nil
	}

	id := inferSpecialIDFromPath(path)
	if id == "" {
		return nil
	}
	moduleRoot, err := findModuleRoot()
	if err != nil {
		return err
	}
	attributePath := filepath.Join(moduleRoot, "fighter", "attributes", strings.ToLower(id)+".go")
	if fileExists(attributePath) {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("No special attribute file found for %s at %s.\n", firstRoot, attributePath)
	fmt.Printf("Create it as type %s now? [Y/n]: ", id)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	if input != "" && !strings.EqualFold(input, "y") && !strings.EqualFold(input, "yes") {
		return nil
	}
	return writeSpecialAttributes(path, data, desc, id, fieldType, false)
}

func inferSpecialIDFromPath(path string) string {
	name := filepath.Base(path)
	matches := regexp.MustCompile(`(?i)^Pl([A-Za-z0-9]{2})\.dat$`).FindStringSubmatch(name)
	if len(matches) != 2 {
		return ""
	}
	code := matches[1]
	return strings.ToUpper(code[:1]) + strings.ToLower(code[1:])
}

func generateSpecialAttributeSource(id string, fieldTypes []string) ([]byte, error) {
	var b strings.Builder
	fmt.Fprintln(&b, "package attributes")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, `import "github.com/Gurvan/melee-data-tools/helpers"`)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "type %s struct {\n", id)
	for i, fieldType := range fieldTypes {
		offset := i * 4
		fmt.Fprintf(&b, "\t%s %s // 0x%X\n", unknownFieldName(offset), fieldType, offset)
	}
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "func (a %s) String() string {\n", id)
	fmt.Fprintln(&b, "\treturn helpers.PrettyString(a)")
	fmt.Fprintln(&b, "}")

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		return nil, err
	}
	return formatted, nil
}

func updateSpecialSwitch(path string, rootName string, id string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	source := string(data)
	newCase := fmt.Sprintf("case %q:\n\t\treturn any(&%s{}), nil", rootName, id)
	casePattern := regexp.MustCompile(`case\s+"` + regexp.QuoteMeta(rootName) + `":\s*\n\s*return any\(&[A-Za-z][A-Za-z0-9]*\{\}\), nil`)
	if casePattern.MatchString(source) {
		source = casePattern.ReplaceAllString(source, newCase)
	} else {
		defaultMarker := "\tdefault:"
		if !strings.Contains(source, defaultMarker) {
			return fmt.Errorf("could not find default case in %s", path)
		}
		source = strings.Replace(source, defaultMarker, "\t"+newCase+"\n\tdefault:", 1)
	}

	formatted, err := format.Source([]byte(source))
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0o644)
}

func unknownFieldName(offset int) string {
	if offset < 0x100 {
		return fmt.Sprintf("Unk0x%02X", offset)
	}
	return fmt.Sprintf("Unk0x%X", offset)
}

func validateSpecialID(id string) error {
	if matched := regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`).MatchString(id); !matched {
		return fmt.Errorf("invalid special id %q; use a Go exported identifier like Fc, Pr, or Ca", id)
	}
	return nil
}

func isSupportedFieldType(fieldType string) bool {
	switch fieldType {
	case "auto", "uint32", "int32", "float32", "Hex32":
		return true
	default:
		return false
	}
}

func specialFieldTypes(path string, offset lib.Addr, count int, mode string) ([]string, error) {
	if mode != "auto" {
		fieldTypes := make([]string, count)
		for i := range fieldTypes {
			fieldTypes[i] = mode
		}
		return fieldTypes, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	start := int(offset)
	end := start + count*4
	if start < 0 || end > len(raw) {
		return nil, fmt.Errorf("special attribute range 0x%X..0x%X is outside %s", start, end, path)
	}

	fieldTypes := make([]string, count)
	for i := 0; i < count; i++ {
		bits := binary.BigEndian.Uint32(raw[start+i*4 : start+i*4+4])
		fieldTypes[i] = inferSpecialFieldType(bits)
	}
	return fieldTypes, nil
}

func inferSpecialFieldType(bits uint32) string {
	if bits == 0 {
		return "float32"
	}

	floatValue := math.Float32frombits(bits)
	intValue := int32(bits)
	if suspiciousFloat(floatValue) && cleanSmallInt(intValue) {
		return "int32"
	}
	if suspiciousFloat(floatValue) && cleanBytePattern(bits) {
		return "Hex32"
	}
	return "float32"
}

func suspiciousFloat(value float32) bool {
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return true
	}
	abs := math.Abs(float64(value))
	return abs != 0 && (abs < 1e-20 || abs > 1e10)
}

func cleanSmallInt(value int32) bool {
	return value >= -10000 && value <= 10000
}

func cleanBytePattern(bits uint32) bool {
	// Some flags/sentinels look like impossible floats but have very deliberate
	// byte patterns, such as 0xff000000 or 0xff00ffff.
	bytes := []byte{
		byte(bits >> 24),
		byte(bits >> 16),
		byte(bits >> 8),
		byte(bits),
	}
	for _, b := range bytes {
		if b != 0x00 && b != 0xff {
			return false
		}
	}
	return true
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if fileExists(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod from %s", dir)
		}
		dir = parent
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func renderRelocation(desc descriptor.Descriptor) string {
	offsets := make([]uint32, 0, len(desc.Relocation))
	for offset := range desc.Relocation {
		offsets = append(offsets, uint32(offset))
	}
	sort.Slice(offsets, func(i, j int) bool {
		return offsets[i] < offsets[j]
	})

	var b strings.Builder
	for _, offset := range offsets {
		fmt.Fprintf(&b, "0x%X -> %d\n", offset, desc.Relocation[lib.Addr(offset)])
	}
	return b.String()
}

func blockSize(desc descriptor.Descriptor, offset lib.Addr) uint32 {
	if offset == 0 {
		return 0
	}
	return desc.Relocation[offset]
}

func renderItems(data mdt.FighterData) string {
	if len(data.Items) == 0 {
		return "No items parsed.\n"
	}

	var b strings.Builder
	for i, item := range data.Items {
		fmt.Fprintf(&b, "Item %d\n", i)
		fmt.Fprintf(&b, "  Attributes: %s\n", yesNo(item.Attributes.ValuePtr != nil))
		fmt.Fprintf(&b, "  Hurtboxes: %d\n", len(item.Hurtboxes))
		fmt.Fprintf(&b, "  States: %d\n", len(item.States))
		if item.Model.ValuePtr != nil {
			fmt.Fprintf(&b, "  Model bone count: %d\n", item.Model.ValuePtr.BoneCount)
			fmt.Fprintf(&b, "  Model bone attach id: %d\n", item.Model.ValuePtr.BoneAttachID)
		} else {
			fmt.Fprintf(&b, "  Model: no\n")
		}
		fmt.Fprintln(&b)
	}
	return b.String()
}

func renderShieldPose(data mdt.FighterData) string {
	if data.ShieldPose.ValuePtr == nil {
		return "No shield pose parsed.\n"
	}

	var b strings.Builder
	shieldPose := data.ShieldPose.ValuePtr
	fmt.Fprintf(&b, "Offset: %s\n", data.ShieldPose.Offset.String())
	fmt.Fprintf(&b, "root: %s\n", optionalJointOffset(shieldPose.Root))
	fmt.Fprintf(&b, "unknown: %s\n", shieldPose.Unknown.String())
	fmt.Fprintln(&b)

	if shieldPose.Root.ValuePtr == nil {
		fmt.Fprintln(&b, "No shield pose root joint tree parsed.")
		return b.String()
	}

	fmt.Fprintln(&b, "Root joint tree:")
	renderJointTreeCompact(&b, shieldPose.Root.ValuePtr, shieldPose.Root.Offset, 0, 0)

	poseJoint := shieldPose.Joint()
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Pose joint tree:")
	fmt.Fprintln(&b, "  Reference code uses the root joint's child pointer as the shield pose.")
	fmt.Fprintf(&b, "  pose offset: %s\n", optionalJointOffset(shieldPose.Root.ValuePtr.Child))
	if poseJoint != nil {
		renderJointTreeCompact(&b, poseJoint, shieldPose.Root.ValuePtr.Child.Offset, 0, 0)
	}

	return b.String()
}

func optionalJointOffset(ptr lib.OptionalPtr[model.Joint]) string {
	if ptr.ValuePtr == nil {
		return "<nil>"
	}
	return ptr.Offset.String()
}

func renderJointTreeCompact(b *strings.Builder, joint *model.Joint, offset lib.Addr, depth int, index int) int {
	if joint == nil {
		return index
	}
	fmt.Fprintf(
		b,
		"%03d %s%s  rot:(% .6f,% .6f,% .6f)  trans:(% .6f,% .6f,% .6f)  child:%s  next:%s\n",
		index,
		strings.Repeat("  ", depth),
		offset.String(),
		joint.Rotation.X,
		joint.Rotation.Y,
		joint.Rotation.Z,
		joint.Translation.X,
		joint.Translation.Y,
		joint.Translation.Z,
		optionalJointOffset(joint.Child),
		optionalJointOffset(joint.Next),
	)
	index++
	if joint.Child.ValuePtr != nil {
		index = renderJointTreeCompact(b, joint.Child.ValuePtr, joint.Child.Offset, depth+1, index)
	}
	if joint.Next.ValuePtr != nil {
		index = renderJointTreeCompact(b, joint.Next.ValuePtr, joint.Next.Offset, depth, index)
	}
	return index
}

func renderActions(data mdt.FighterData) string {
	if data.ActionTable.ValuePtr == nil || len(*data.ActionTable.ValuePtr) == 0 {
		return "No actions parsed.\n"
	}

	var b strings.Builder
	for i, action := range *data.ActionTable.ValuePtr {
		fmt.Fprintf(&b, "%03d  %-40s subactions:%3d  anim_size:%d  flags:%+v\n",
			i,
			actionName(action),
			len(action.Subactions),
			action.AnimationSize,
			action.Flags,
		)
	}
	return b.String()
}

func renderAnimations(data mdt.FighterData, animations []reviewAnimation) string {
	actionRefs := animationActionRefs(data)
	if len(animations) == 0 {
		if len(actionRefs) == 0 {
			return "No animation bundle parsed and no action animation references found.\nRun with -animations /path/to/PlXXAJ.dat to list every animation in the animation bundle.\n"
		}
		return renderActionAnimationRefs(actionRefs)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Animations: %d\n", len(animations))
	fmt.Fprintf(&b, "Action animation references: %d\n\n", len(actionRefs))
	for _, animation := range animations {
		refs := actionRefs[animation.Offset]
		fmt.Fprintf(&b, "%03d  offset:0x%X  size:%d  root:%s  frames:%g  nodes:%d  tracks:%d  refs:%d\n",
			animation.Index,
			animation.Offset,
			animation.Size,
			emptyDefault(animation.Root, "<unknown>"),
			animation.Data.FrameCount,
			len(animation.Data.TracksCounts),
			animationTrackCount(animation.Data),
			len(refs),
		)
		if len(refs) > 0 {
			for _, ref := range refs {
				fmt.Fprintf(&b, "      action %03d  %-40s anim_size:%d\n", ref.Index, ref.Name, ref.Size)
			}
		}
	}

	var unparsed []int64
	for offset := range actionRefs {
		if !hasAnimationOffset(animations, offset) {
			unparsed = append(unparsed, offset)
		}
	}
	sort.Slice(unparsed, func(i, j int) bool { return unparsed[i] < unparsed[j] })
	if len(unparsed) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "Action references missing from the parsed animation bundle:")
		for _, offset := range unparsed {
			fmt.Fprintf(&b, "  offset:0x%X\n", offset)
			for _, ref := range actionRefs[offset] {
				fmt.Fprintf(&b, "      action %03d  %-40s anim_size:%d\n", ref.Index, ref.Name, ref.Size)
			}
		}
	}

	return b.String()
}

func renderActionAnimationRefs(actionRefs map[int64][]actionAnimationRef) string {
	offsets := make([]int64, 0, len(actionRefs))
	for offset := range actionRefs {
		offsets = append(offsets, offset)
	}
	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })

	var b strings.Builder
	fmt.Fprintln(&b, "No animation bundle parsed.")
	fmt.Fprintln(&b, "Showing only animation offsets referenced by actions.")
	fmt.Fprintln(&b, "Run with -animations /path/to/PlXXAJ.dat to list every animation in the animation bundle.")
	fmt.Fprintln(&b)
	for _, offset := range offsets {
		refs := actionRefs[offset]
		fmt.Fprintf(&b, "offset:0x%X  refs:%d\n", offset, len(refs))
		for _, ref := range refs {
			fmt.Fprintf(&b, "    action %03d  %-40s anim_size:%d\n", ref.Index, ref.Name, ref.Size)
		}
	}
	return b.String()
}

func renderSubactions(data mdt.FighterData) string {
	if data.ActionTable.ValuePtr == nil || len(*data.ActionTable.ValuePtr) == 0 {
		return "No actions parsed.\n"
	}

	var b strings.Builder
	for i, action := range *data.ActionTable.ValuePtr {
		if len(action.Subactions) == 0 {
			continue
		}
		fmt.Fprintf(&b, "%03d %s\n", i, actionName(action))
		for j, subaction := range action.Subactions {
			fmt.Fprintf(&b, "  %03d  %T  %+v\n", j, subaction, subaction)
		}
		fmt.Fprintln(&b)
	}
	if b.Len() == 0 {
		return "Actions parsed, but none had subactions.\n"
	}
	return b.String()
}

func readAnimationBundle(path string) ([]reviewAnimation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	files, offsets, err := mdt.SplitAnimationFile(binread.NewReader(file))
	if err != nil {
		return nil, err
	}

	animations := make([]reviewAnimation, 0, len(files))
	for i, fileData := range files {
		var animationFile mdt.AnimationFile
		animationData, desc, err := animationFile.ReadFromBytes(fileData)
		if err != nil {
			return nil, fmt.Errorf("animation %d at 0x%X: %w", i, offsets[i], err)
		}
		root, err := desc.FirstRootName()
		if err != nil {
			root = ""
		}
		animations = append(animations, reviewAnimation{
			Index:  i,
			Offset: offsets[i],
			Size:   len(fileData),
			Root:   root,
			Data:   animationData,
		})
	}
	return animations, nil
}

func animationActionRefs(data mdt.FighterData) map[int64][]actionAnimationRef {
	refs := make(map[int64][]actionAnimationRef)
	if data.ActionTable.ValuePtr == nil {
		return refs
	}
	for i, action := range *data.ActionTable.ValuePtr {
		if action.AnimationSize == 0 || action.Animation.Offset == 0 {
			continue
		}
		offset := action.Animation.Offset.ToSeek() - 0x20
		refs[offset] = append(refs[offset], actionAnimationRef{
			Index: i,
			Name:  actionName(action),
			Size:  action.AnimationSize,
		})
	}
	return refs
}

func animationTrackCount(animation mdt.AnimationData) int {
	if animation.Tracks.ValuePtr == nil {
		return 0
	}
	count := 0
	for _, tracks := range *animation.Tracks.ValuePtr {
		count += len(tracks)
	}
	return count
}

func hasAnimationOffset(animations []reviewAnimation, offset int64) bool {
	for _, animation := range animations {
		if animation.Offset == offset {
			return true
		}
	}
	return false
}

func emptyDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func runMenu(sections []section) {
	reader := bufio.NewReader(os.Stdin)
	for {
		clearScreen()
		fmt.Println("Fighter Review")
		fmt.Println(strings.Repeat("=", 14))
		for i, section := range sections {
			fmt.Printf("%2d. %-16s %s\n", i+1, section.key, section.title)
		}
		fmt.Println()
		fmt.Print("Choose a section number/key, or q to quit: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if strings.EqualFold(input, "q") || strings.EqualFold(input, "quit") {
			return
		}

		section, ok := selectSection(sections, input)
		if !ok {
			fmt.Printf("\nUnknown section %q. Press Enter to continue.", input)
			_, _ = reader.ReadString('\n')
			continue
		}

		clearScreen()
		printSection(section)
		fmt.Print("\nPress Enter to return to the menu.")
		_, _ = reader.ReadString('\n')
	}
}

func selectSection(sections []section, input string) (section, bool) {
	if idx, err := strconv.Atoi(input); err == nil {
		if idx >= 1 && idx <= len(sections) {
			return sections[idx-1], true
		}
		return section{}, false
	}
	return findSection(sections, input)
}

func findSection(sections []section, key string) (section, bool) {
	key = strings.ToLower(strings.TrimSpace(key))
	for _, section := range sections {
		if section.key == key {
			return section, true
		}
	}
	return section{}, false
}

func printSection(section section) {
	fmt.Println(section.title)
	fmt.Println(strings.Repeat("=", len(section.title)))
	fmt.Println(section.body())
}

func dumpPtr(value interface{}) string {
	if value == nil {
		return "No data parsed.\n"
	}
	return dump(value)
}

func dump(value interface{}) string {
	config := spew.ConfigState{
		DisableCapacities:       true,
		DisablePointerAddresses: true,
		Indent:                  "  ",
	}
	return config.Sdump(value)
}

func actionName(action fighter.Action) string {
	if action.Name.ValuePtr == nil {
		return "<unnamed>"
	}
	name := action.Name.ValuePtr.String()
	if name == "" {
		return "<unnamed>"
	}
	return name
}

func lenActionTable(actions *fighter.ActionTable) int {
	if actions == nil {
		return 0
	}
	return len(*actions)
}

func lenHurtboxes(hurtboxes *fighter.Hurtboxes) int {
	if hurtboxes == nil {
		return 0
	}
	return len(*hurtboxes)
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
