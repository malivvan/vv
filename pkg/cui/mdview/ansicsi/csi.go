package ansicsi

import (
	"bytes"
	"io"
	"strconv"
)

// Command represents a parsed ANSI control function.
type Command interface {
	// Encode writes the ANSI CSI and control sequence for the command to the given Writer.
	Encode(w io.Writer) (int, error)

	decodeParameters(params []int) bool
}

// ControlSequence represents a single ANSI control sequence.
type ControlSequence struct {
	Parameters   []byte
	Intermediate []byte
	Final        byte
}

func (cs *ControlSequence) Encode(w io.Writer) (int, error) {
	bytes := make([]byte, 0, 2+len(cs.Parameters)+len(cs.Intermediate)+1)
	bytes = append(bytes, []byte("\x1b[")...)
	bytes = append(bytes, cs.Parameters...)
	bytes = append(bytes, cs.Intermediate...)
	bytes = append(bytes, cs.Final)
	return w.Write(bytes)
}

func (*ControlSequence) decodeParameters(params []int) bool {
	return false
}

// Decode decodes the ANSI control function beginning at the first byte of b and returns the function, its
// parameters, and its encoded size. If a valid control sequence is found but the control function is not
// recognized, the raw control sequence is returned as a *ControlSequence value.
func Decode(b []byte) (Command, int) {
	if len(b) < 2 || b[0] != 0x1b || b[1] != '[' {
		return nil, 0
	}
	b = b[2:]

	// parameter bytes
	paramEnd := 0
	for paramEnd < len(b) && b[paramEnd] >= 0x30 && b[paramEnd] < 0x40 {
		paramEnd++
	}
	params, b := b[:paramEnd], b[paramEnd:]

	// intermediate bytes
	intermediateEnd := 0
	for intermediateEnd < len(b) && b[intermediateEnd] >= 0x20 && b[intermediateEnd] < 0x30 {
		intermediateEnd++
	}
	intermediate, b := b[:intermediateEnd], b[intermediateEnd:]

	// final byte
	if len(b) < 1 || b[0] < 0x40 || b[0] > 0x7e {
		return nil, 0
	}
	final := b[0]

	size := 2 + len(params) + len(intermediate) + 1
	cmd, ok := decodeCommand(params, intermediate, final)
	if !ok {
		cmd = &ControlSequence{
			Parameters:   params,
			Intermediate: intermediate,
			Final:        final,
		}
	}
	return cmd, size
}

func getCommand(intermediate []byte, final byte) (Command, bool) {
	switch {
	case len(intermediate) == 0:
		switch final {
		case 0x40: // 8.3.64 ICH - INSERT CHARACTER (Pn)
			return nil, false
		case 0x41: // 8.3.22 CUU - CURSOR UP (Pn)
			return nil, false
		case 0x42: // 8.3.19 CUD - CURSOR DOWN (Pn)
			return nil, false
		case 0x43: // 8.3.20 CUF - CURSOR RIGHT (Pn)
			return nil, false
		case 0x44: // 8.3.18 CUB - CURSOR LEFT (Pn)
			return nil, false
		case 0x45: // 8.3.12 CNL - CURSOR NEXT LINE (Pn)
			return nil, false
		case 0x46: // 8.3.13 CPL - CURSOR PRECEDING LINE (Pn)
			return nil, false
		case 0x47: // 8.3.9 CHA - CURSOR CHARACTER ABSOLUTE (Pn)
			return nil, false
		case 0x48: // 8.3.21 CUP - CURSOR POSITION (Pn1;Pn2)
			return nil, false
		case 0x49: // 8.3.10 CHT - CURSOR FORWARD TABULATION (Pn)
			return nil, false
		case 0x4a: // 8.3.39 ED - ERASE IN PAGE (Ps)
			return nil, false
		case 0x4b: // 8.3.41 EL - ERASE IN LINE (Ps)
			return nil, false
		case 0x4c: // 8.3.67 IL - INSERT LINE (Pn)
			return nil, false
		case 0x4d: // 8.3.32 DL - DELETE LINE (Pn)
			return nil, false
		case 0x4e: // 8.3.40 EF - ERASE IN FIELD (Ps)
			return nil, false
		case 0x4f: // 8.3.37 EA - ERASE IN AREA (Ps)
			return nil, false
		case 0x50: // 8.3.26 DCH - DELETE CHARACTER (Pn)
			return nil, false
		case 0x51: // 8.3.115 SEE - SELECT EDITING EXTENT (Ps)
			return nil, false
		case 0x52: // 8.3.14 CPR - ACTIVE POSITION REPORT (Pn1;Pn2)
			return nil, false
		case 0x53: // 8.3.147 SU - SCROLL UP (Pn)
			return nil, false
		case 0x54: // 8.3.113 SD - SCROLL DOWN (Pn)
			return nil, false
		case 0x55: // 8.3.87 NP - NEXT PAGE (Pn)
			return nil, false
		case 0x56: // 8.3.95 PP - PRECEDING PAGE (Pn)
			return nil, false
		case 0x57: // 8.3.17 CTC - CURSOR TABULATION CONTROL (Ps...)
			return nil, false
		case 0x58: // 8.3.38 ECH - ERASE CHARACTER (Pn)
			return nil, false
		case 0x59: // 8.3.23 CVT - CURSOR LINE TABULATION (Pn)
			return nil, false
		case 0x5a: // 8.3.7 CBT - CURSOR BACKWARD TABULATION (Pn)
			return nil, false
		case 0x5b: // 8.3.137 SRS - START REVERSED STRING (Ps)
			return nil, false
		case 0x5c: // 8.3.99 PTX - PARALLEL TEXTS (Ps)
			return nil, false
		case 0x5d: // 8.3.114 SDS - START DIRECTED STRING (Ps)
			return nil, false
		case 0x5e: // 8.3.120 SIMD - SELECT IMPLICIT MOVEMENT DIRECTION (Ps)
			return nil, false
		case 0x60: // 8.3.57 HPA - CHARACTER POSITION ABSOLUTE (Pn)
			return nil, false
		case 0x61: // 8.3.59 HPR - CHARACTER POSITION FORWARD (Pn)
			return nil, false
		case 0x62: // 8.3.103 REP - REPEAT (Pn)
			return nil, false
		case 0x63: // 8.3.24 DA - DEVICE ATTRIBUTES (Ps)
			return nil, false
		case 0x64: // 8.3.158 VPA - LINE POSITION ABSOLUTE (Pn)
			return nil, false
		case 0x65: // 8.3.160 VPR - LINE POSITION FORWARD (Pn)
			return nil, false
		case 0x66: // 8.3.63 HVP - CHARACTER AND LINE POSITION (Pn1;Pn2)
			return nil, false
		case 0x67: // 8.3.154 TBC - TABULATION CLEAR (Ps)
			return nil, false
		case 0x68: // 8.3.125 SM - SET MODE (Ps...)
			return nil, false
		case 0x69: // 8.3.82 MC - MEDIA COPY (Ps)
			return nil, false
		case 0x6a: // 8.3.58 HPB - CHARACTER POSITION BACKWARD (Pn)
			return nil, false
		case 0x6b: // 8.3.159 VPB - LINE POSITION BACKWARD (Pn)
			return nil, false
		case 0x6c: // 8.3.106 RM - RESET MODE (Ps...)
			return nil, false
		case 0x6d: // 8.3.117 SGR - SELECT GRAPHIC RENDITION (Ps...)
			return &SetGraphicsRendition{}, true
		case 0x6e: // 8.3.35 DSR - DEVICE STATUS REPORT (Ps)
			return nil, false
		case 0x6f: // 8.3.25 DAQ - DEFINE AREA QUALIFICATION (Ps...)
			return nil, false
		}
	case len(intermediate) == 1 && intermediate[0] == 0x20:
		switch final {
		case 0x40: // 8.3.121 SL - SCROLL LEFT (Pn)
			return nil, false
		case 0x41: // 8.3.135 SR - SCROLL RIGHT (Pn)
			return nil, false
		case 0x42: // 8.3.55 GSM - GRAPHIC SIZE MODIFICATION (Pn1;Pn2)
			return nil, false
		case 0x43: // 8.3.56 GSS - GRAPHIC SIZE SELECTION (Pn)
			return nil, false
		case 0x44: // 8.3.53 FNT - FONT SELECTION (Ps1;Ps2)
			return nil, false
		case 0x45: // 8.3.157 TSS - THIN SPACE SPECIFICATION (Pn)
			return nil, false
		case 0x46: // 8.3.73 JFY - JUSTIFY (Ps...)
			return nil, false
		case 0x47: // 8.3.132 SPI - SPACING INCREMENT (Pn1;Pn2)
			return nil, false
		case 0x48: // 8.3.102 QUAD - QUAD (Ps...)
			return nil, false
		case 0x49: // 8.3.139 SSU - SELECT SIZE UNIT (Ps)
			return nil, false
		case 0x4a: // 8.3.91 PFS - PAGE FORMAT SELECTION (Ps)
			return nil, false
		case 0x4b: // 8.3.118 SHS - SELECT CHARACTER SPACING (Ps)
			return nil, false
		case 0x4c: // 8.3.149 SVS - SELECT LINE SPACING (Ps)
			return nil, false
		case 0x4d: // 8.3.66 IGS - IDENTIFY GRAPHIC SUBREPERTOIRE (Ps)
			return nil, false
		case 0x4f: // 8.3.65 IDCS - IDENTIFY DEVICE CONTROL STRING (Ps)
			return nil, false
		case 0x50: // 8.3.96 PPA - PAGE POSITION ABSOLUTE (Pn)
			return nil, false
		case 0x51: // 8.3.98 PPR - PAGE POSITION FORWARD (Pn)
			return nil, false
		case 0x52: // 8.3.97 PPB - PAGE POSITION BACKWARD (Pn)
			return nil, false
		case 0x53: // 8.3.130 SPD - SELECT PRESENTATION DIRECTIONS (Ps1;Ps2)
			return nil, false
		case 0x54: // 8.3.36 DTA - DIMENSION TEXT AREA (Pn1;Pn2)
			return nil, false
		case 0x55: // 8.3.122 SLH - SET LINE HOME (Pn)
			return nil, false
		case 0x56: // 8.3.123 SLL - SET LINE LIMIT (Pn)
			return nil, false
		case 0x57: // 8.3.52 FNK - FUNCTION KEY (Pn)
			return nil, false
		case 0x58: // 8.3.134 SPQR - SELECT PRINT QUALITY AND RAPIDITY (Ps)
			return nil, false
		case 0x59: // 8.3.116 SEF - SHEET EJECT AND FEED (Ps1;Ps2)
			return nil, false
		case 0x5a: // 8.3.90 PEC - PRESENTATION EXPAND OR CONTRACT (Ps)
			return nil, false
		case 0x5b: // 8.3.140 SSW - SET SPACE WIDTH (Pn)
			return nil, false
		case 0x5c: // 8.3.107 SACS - SET ADDITIONAL CHARACTER SEPARATION (Pn)
			return nil, false
		case 0x5d: // 8.3.108 SAPV - SELECT ALTERNATIVE PRESENTATION VARIANTS (Ps...)
			return nil, false
		case 0x5e: // 8.3.144 STAB - SELECTIVE TABULATION (Ps)
			return nil, false
		case 0x5f: // 8.3.54 GCC - GRAPHIC CHARACTER COMBINATION (Ps)
			return nil, false
		case 0x60: // 8.3.153 TATE - TABULATION ALIGNED TRAILING EDGE (Pn)
			return nil, false
		case 0x61: // 8.3.152 TALE - TABULATION ALIGNED LEADING EDGE (Pn)
			return nil, false
		case 0x62: // 8.3.151 TAC - TABULATION ALIGNED CENTRED (Pn)
			return nil, false
		case 0x63: // 8.3.155 TCC - TABULATION CENTRED ON CHARACTER (Pn1;Pn2)
			return nil, false
		case 0x64: // 8.3.156 TSR - TABULATION STOP REMOVE (Pn)
			return nil, false
		case 0x65: // 8.3.110 SCO - SELECT CHARACTER ORIENTATION (Ps)
			return nil, false
		case 0x66: // 8.3.136 SRCS - SET REDUCED CHARACTER SEPARATION (Pn)
			return nil, false
		case 0x67: // 8.3.112 SCS - SET CHARACTER SPACING (Pn)
			return nil, false
		case 0x68: // 8.3.124 SLS - SET LINE SPACING (Pn)
			return nil, false
		case 0x69: // 8.3.131 SPH - SET PAGE HOME (Pn)
			return nil, false
		case 0x6a: // 8.3.133 SPL - SET PAGE LIMIT (Pn)
			return nil, false
		case 0x6b: // 8.3.111 SCP - SELECT CHARACTER PATH (Ps1;Ps2)
			return nil, false
		}
	}
	return nil, false
}

func decodeCommand(parameters, intermediate []byte, final byte) (Command, bool) {
	cmd, ok := getCommand(intermediate, final)
	if !ok {
		return nil, false
	}

	// Decode the parameter list.
	var param []byte
	var params []int
	for rawParams := parameters; len(rawParams) > 0; rawParams = rawParams[1:] {
		if rawParams[0] == ';' {
			if len(param) == 0 {
				params = append(params, -1)
			} else {
				i, err := strconv.ParseUint(string(param), 10, 0)
				if err != nil {
					return nil, false
				}
				param, params = param[:0], append(params, int(i))
			}
		} else {
			param = append(param, rawParams[0])
		}
	}
	if len(param) != 0 {
		i, err := strconv.ParseUint(string(param), 10, 0)
		if err != nil {
			return nil, false
		}
		param, params = param[:0], append(params, int(i))
	}
	if len(parameters) != 0 && parameters[len(parameters)-1] == ';' {
		params = append(params, -1)
	}

	return cmd, cmd.decodeParameters(params)
}

func encodeCommand(w io.Writer, parameters []int, intermediate []byte, final byte) (int, error) {
	// Encode the parameter list.
	var params bytes.Buffer
	for i, p := range parameters {
		if i > 0 {
			params.WriteByte(';')
		}
		if p >= 0 {
			params.WriteString(strconv.FormatUint(uint64(p), 10))
		}
	}

	// Write the CSI + control sequence.
	cs := ControlSequence{
		Parameters:   params.Bytes(),
		Intermediate: intermediate,
		Final:        final,
	}
	return cs.Encode(w)
}
