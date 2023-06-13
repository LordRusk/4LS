package conn

import "github.com/lordrusk/4LS/bord"

/* BIG so it gets its own file */
type allbrds struct {
	pointerMap map[string]*Bored

	A, B, C, D, E, F, G, GIF, H, HR, K, M, O, P, R, S, T, U, V,
	VG, VM, VMG, VR, VRPG, VST, W, WG, I, IC, R9K, S4S, VIP, QA,
	CM, HM, LGBT, Y, THR, ACO, ADV, AN, BANT, BIZ, CGL, CK, CO,
	DIY, FA, FIT, GD, HC, HIS, INT, JP, LIT, MLP, MU, N, NEWS,
	OUT, PO, POL, PW, QST, SCI, SOC, SP, TG, TOY, TRV, TV, VP,
	VT, WSG, WSR, X, XS Bored
}

/* essentially legacy conn.Con structure support */
func (a *allbrds) All(board string) *Bored {
	return a.pointerMap[board]

}

/* essentially legacy conn.Con structure support */
func (a *allbrds) Map() map[string]*Bored {
	return a.pointerMap
}

func mkAll() *allbrds {
	/* extremely helpful */
	a := allbrds{
		A:    mkBored(bord.A),
		B:    mkBored(bord.B),
		C:    mkBored(bord.C),
		D:    mkBored(bord.D),
		E:    mkBored(bord.E),
		F:    mkBored(bord.F),
		G:    mkBored(bord.G),
		GIF:  mkBored(bord.GIF),
		H:    mkBored(bord.H),
		HR:   mkBored(bord.HR),
		K:    mkBored(bord.K),
		M:    mkBored(bord.M),
		O:    mkBored(bord.O),
		P:    mkBored(bord.P),
		R:    mkBored(bord.R),
		S:    mkBored(bord.S),
		T:    mkBored(bord.T),
		U:    mkBored(bord.U),
		V:    mkBored(bord.V),
		VG:   mkBored(bord.VG),
		VM:   mkBored(bord.VM),
		VMG:  mkBored(bord.VMG),
		VR:   mkBored(bord.VR),
		VRPG: mkBored(bord.VRPG),
		VST:  mkBored(bord.VST),
		W:    mkBored(bord.W),
		WG:   mkBored(bord.WG),
		I:    mkBored(bord.I),
		IC:   mkBored(bord.IC),
		R9K:  mkBored(bord.R9K),
		S4S:  mkBored(bord.S4S),
		VIP:  mkBored(bord.VIP),
		QA:   mkBored(bord.QA),
		CM:   mkBored(bord.CM),
		HM:   mkBored(bord.HM),
		LGBT: mkBored(bord.LGBT),
		Y:    mkBored(bord.Y),
		THR:  mkBored(bord.THR),
		ACO:  mkBored(bord.ACO),
		ADV:  mkBored(bord.ADV),
		AN:   mkBored(bord.AN),
		BANT: mkBored(bord.BANT),
		BIZ:  mkBored(bord.BIZ),
		CGL:  mkBored(bord.CGL),
		CK:   mkBored(bord.CK),
		CO:   mkBored(bord.CO),
		DIY:  mkBored(bord.DIY),
		FA:   mkBored(bord.FA),
		FIT:  mkBored(bord.FIT),
		GD:   mkBored(bord.GD),
		HC:   mkBored(bord.HC),
		HIS:  mkBored(bord.HIS),
		INT:  mkBored(bord.INT),
		JP:   mkBored(bord.JP),
		LIT:  mkBored(bord.LIT),
		MLP:  mkBored(bord.MLP),
		MU:   mkBored(bord.MU),
		N:    mkBored(bord.N),
		NEWS: mkBored(bord.NEWS),
		OUT:  mkBored(bord.OUT),
		PO:   mkBored(bord.PO),
		POL:  mkBored(bord.POL),
		PW:   mkBored(bord.PW),
		QST:  mkBored(bord.QST),
		SCI:  mkBored(bord.SCI),
		SOC:  mkBored(bord.SOC),
		SP:   mkBored(bord.SP),
		TG:   mkBored(bord.TG),
		TOY:  mkBored(bord.TOY),
		TRV:  mkBored(bord.TRV),
		TV:   mkBored(bord.TV),
		VP:   mkBored(bord.VP),
		VT:   mkBored(bord.VT),
		WSG:  mkBored(bord.WSG),
		WSR:  mkBored(bord.WSR),
		X:    mkBored(bord.X),
		XS:   mkBored(bord.XS),
	}
	m := make(map[string]*Bored)
	m["A"] = &a.A
	m["B"] = &a.B
	m["C"] = &a.C
	m["D"] = &a.D
	m["E"] = &a.E
	m["F"] = &a.F
	m["G"] = &a.G
	m["GIF"] = &a.GIF
	m["H"] = &a.H
	m["HR"] = &a.HR
	m["K"] = &a.K
	m["M"] = &a.M
	m["O"] = &a.O
	m["P"] = &a.P
	m["R"] = &a.R
	m["S"] = &a.S
	m["T"] = &a.T
	m["U"] = &a.U
	m["V"] = &a.V
	m["VG"] = &a.VG
	m["VM"] = &a.VM
	m["VMG"] = &a.VMG
	m["VR"] = &a.VR
	m["VRPG"] = &a.VRPG
	m["VST"] = &a.VST
	m["W"] = &a.W
	m["WG"] = &a.WG
	m["I"] = &a.I
	m["IC"] = &a.IC
	m["R9K"] = &a.R9K
	m["S4S"] = &a.S4S
	m["VIP"] = &a.VIP
	m["QA"] = &a.QA
	m["CM"] = &a.CM
	m["HM"] = &a.HM
	m["LGBT"] = &a.LGBT
	m["Y"] = &a.Y
	m["THR"] = &a.THR
	m["ACO"] = &a.ACO
	m["ADV"] = &a.ADV
	m["AN"] = &a.AN
	m["BANT"] = &a.BANT
	m["BIZ"] = &a.BIZ
	m["CGL"] = &a.CGL
	m["CK"] = &a.CK
	m["CO"] = &a.CO
	m["DIY"] = &a.DIY
	m["FA"] = &a.FA
	m["FIT"] = &a.FIT
	m["GD"] = &a.GD
	m["HC"] = &a.HC
	m["HIS"] = &a.HIS
	m["INT"] = &a.INT
	m["JP"] = &a.JP
	m["LIT"] = &a.LIT
	m["MLP"] = &a.MLP
	m["MU"] = &a.MU
	m["N"] = &a.N
	m["NEWS"] = &a.NEWS
	m["OUT"] = &a.OUT
	m["PO"] = &a.PO
	m["POL"] = &a.POL
	m["PW"] = &a.PW
	m["QST"] = &a.QST
	m["SCI"] = &a.SCI
	m["SOC"] = &a.SOC
	m["SP"] = &a.SP
	m["TG"] = &a.TG
	m["TOY"] = &a.TOY
	m["TRV"] = &a.TRV
	m["TV"] = &a.TV
	m["VP"] = &a.VP
	m["VT"] = &a.VT
	m["WSG"] = &a.WSG
	m["WSR"] = &a.WSR
	m["X"] = &a.X
	a.pointerMap = m
	return &a
}
