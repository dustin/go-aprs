package aprs

import (
	"bufio"
	"io"
	"log"
	"strings"
)

var primarySymbolMap map[byte]string
var alternateSymbolMap map[byte]string

// Pasted from http://wa8lmf.net/miscinfo/APRSsymbolcodes.txt
var symbolTableTextFile = `
REM Revised by WA8LMF   28 Sept 2005
REM Based on original file by G4IDE supplied with UI-View

REM                Primary Table                  Alternate Table
REM   Symbol   GPSxyz Index    Description     GPSxyz Index    Description
REM   ------   ------ -----    -----------     ------ -----    -----------
        !,      BB,     0,      Police Stn,     OB,     0,      Emergency,
        ",      BC,     1,      No Symbol,      OC,     1,      No Symbol
        #,      BD,     2,      Digi,           OD,     2,      No. Digi
        $,      BE,     3,      Phone,          OE,     3,      Bank, 📱,
        %,      BF,     4,      DX Cluster,     OF,     4,      No Symbol
        &,      BG,     5,      HF Gateway,     OG,     5,      No. Diam'd, Ⓖ,
        ',      BH,     6,      Plane sm,       OH,     6,      Crash site, ✈,
        (,      BI,     7,      Mob Sat Stn,    OI,     7,      Cloudy
        ),      BJ,     8,      WheelChair,     OJ,     8,      MEO, ♿,
        *,      BK,     9,      Snowmobile,     OK,     9,      Snow
        +,      BL,     10,     Red Cross,      OL,     10,     Church
        ,,      BM,     11,     Boy Scout,      OM,     11,     Girl Scout
        -,      BN,     12,     Home,           ON,     12,     Home (HF), ⌂,
        .,      BO,     13,     X,              OO,     13,     UnknownPos
        /,      BP,     14,     Red Dot,        OP,     14,     Destination, ·,
        0,      P0,     15,     Circle (0),     A0,     15,     No. Circle, ⓪,
        1,      P1,     16,     Circle (1),     A1,     16,     No Symbol, ①,
        2,      P2,     17,     Circle (2),     A2,     17,     No Symbol, ②,
        3,      P3,     18,     Circle (3),     A3,     18,     No Symbol, ③,
        4,      P4,     19,     Circle (4),     A4,     19,     No Symbol, ④,
        5,      P5,     20,     Circle (5),     A5,     20,     No Symbol, ⑤,
        6,      P6,     21,     Circle (6),     A6,     21,     No Symbol, ⑥,
        7,      P7,     22,     Circle (7),     A7,     22,     No Symbol, ⑦,
        8,      P8,     23,     Circle (8),     A8,     23,     No Symbol, ⑧,
        9,      P9,     24,     Circle (9),     A9,     24,     Petrol Stn, ⑨,
        :,      MR,     25,     Fire,           NR,     25,     Hail, 🔥,
        ;,      MS,     26,     Campground,     NS,     26,     Park, 🏕,🏞
        <,      MT,     27,     Motorcycle,     NT,     27,     Gale Fl, 🛵,
        =,      MU,     28,     Rail Eng.,      NU,     28,     No Symbol, 🚆,
        >,      MV,     29,     Car,            NV,     29,     No. Car, 🚗,
        ?,      MW,     30,     File svr,       NW,     30,     Info Kiosk
        @,      MX,     31,     HC Future,      NX,     31,     Hurricane
        A,      PA,     32,     Aid Stn,        AA,     32,     No. Box
        B,      PB,     33,     BBS,            AB,     33,     Snow blwng
        C,      PC,     34,     Canoe,          AC,     34,     Coast G'rd, 🛶,
        D,      PD,     35,     No Symbol,      AD,     35,     Drizzle
        E,      PE,     36,     Eyeball,        AE,     36,     Smoke, 👁,
        F,      PF,     37,     Tractor,      AF,     37,     Fr'ze Rain, 🚜,
        G,      PG,     38,     Grid Squ.,      AG,     38,     Snow Shwr
        H,      PH,     39,     Hotel,          AH,     39,     Haze, 🏨,
        I,      PI,     40,     Tcp/ip,         AI,     40,     Rain Shwr
        J,      PJ,     41,     No Symbol,      AJ,     41,     Lightning, , ☇
        K,      PK,     42,     School,         AK,     42,     Kenwood, 🏫,
        L,      PL,     43,     Usr Log-ON,     AL,     43,     Lighthouse, , ⛯
        M,      PM,     44,     MacAPRS,        AM,     44,     No Symbol
        N,      PN,     45,     NTS Stn,        AN,     45,     Nav Buoy
        O,      PO,     46,     Balloon,        AO,     46,     Rocket, 🎈, 🚀
        P,      PP,     47,     Police,         AP,     47,     Parking, 👮,
        Q,      PQ,     48,     TBD,            AQ,     48,     Quake
        R,      PR,     49,     Rec Veh'le,     AR,     49,     Restaurant, 🚙,
        S,      PS,     50,     Shuttle,        AS,     50,     Sat/Pacsat
        T,      PT,     51,     SSTV,           AT,     51,     T'storm
        U,      PU,     52,     Bus,            AU,     52,     Sunny, 🚌,
        V,      PV,     53,     ATV,            AV,     53,     VORTAC
        W,      PW,     54,     WX Service,     AW,     54,     No. WXS
        X,      PX,     55,     Helo,           AX,     55,     Pharmacy
        Y,      PY,     56,     Yacht,          AY,     56,     No Symbol
        Z,      PZ,     57,     WinAPRS,        AZ,     57,     No Symbol
        [,      HS,     58,     Jogger,         DS,     58,     Wall Cloud
        \,      HT,     59,     Triangle,       DT,     59,     No Symbol
        ],      HU,     60,     PBBS,           DU,     60,     No Symbol
        ^,      HV,     61,     Plane lrge,     DV,     61,     No. Plane
        _,      HW,     62,     WX Station,     DW,     62,     No. WX Stn, ☀,
        ` + "`" + `,      HX,     63,     Dish Ant.,      DX,     63,     Rain
        a,      LA,     64,     Ambulance,      SA,     64,     No. Diamond, ☠,
        b,      LB,     65,     Bike,           SB,     65,     Dust blwng, 🚲,
        c,      LC,     66,     ICP,            SC,     66,     No. CivDef
        d,      LD,     67,     Fire Station,   SD,     67,     DX Spot
        e,      LE,     68,     Horse,          SE,     68,     Sleet, 🐎,
        f,      LF,     69,     Fire Truck,     SF,     69,     Funnel Cld, 🚒,
        g,      LG,     70,     Glider,         SG,     70,     Gale
        h,      LH,     71,     Hospital,       SH,     71,     HAM store, 🏥,
        i,      LI,     72,     IOTA,           SI,     72,     No. Blk Box
        j,      LJ,     73,     Jeep,           SJ,     73,     WorkZone
        k,      LK,     74,     Truck,          SK,     74,     SUV
        l,      LL,     75,     Laptop,         SL,     75,     Area Locns
        m,      LM,     76,     Mic-E Rptr,     SM,     76,     Milepost
        n,      LN,     77,     Node,           SN,     77,     No. Triang
        o,      LO,     78,     EOC,            SO,     78,     Circle sm
        p,      LP,     79,     Rover,          SP,     79,     Part Cloud
        q,      LQ,     80,     Grid squ.,      SQ,     80,     No Symbol
        r,      LR,     81,     Antenna,        SR,     81,     Restrooms, ⏉,
        s,      LS,     82,     Power Boat,     SS,     82,     No. Boat
        t,      LT,     83,     Truck Stop,     ST,     83,     Tornado
        u,      LU,     84,     Truck 18wh,     SU,     84,     No. Truck
        v,      LV,     85,     Van,            SV,     85,     No. Van
        w,      LW,     86,     Water Stn,      SW,     86,     Flooding
        x,      LX,     87,     XAPRS,          SX,     87,     No Symbol
        y,      LY,     88,     Yagi,           SY,     88,     Sky Warn, ⏉,
        z,      LZ,     89,     Shelter,        SZ,     89,     No. Shelter
        {,      J1,     90,     No Symbol,      Q1,     90,     Fog
        |,      J2,     91,     TNC Stream Sw,  Q2,     91,     TNC Stream SW
        },      J3,     92,     No Symbol,      Q3,     92,     No Symbol
        ~,      J4,     93,     TNC Stream Sw,  Q4,     93,     TNC Stream SW
`

var symbolGlyphs = map[string]string{}

func init() {
	primarySymbolMap = map[byte]string{}
	alternateSymbolMap = map[byte]string{}

	r := bufio.NewReader(strings.NewReader(symbolTableTextFile))

	for {
		l, err := r.ReadString(byte('\n'))
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Error reading a line:  %v", err)
			}
			break
		}
		l = strings.TrimSpace(l)
		if len(l) < 3 {
			continue
		}
		parts := strings.Split(l[2:], ",")
		if len(parts) >= 6 {
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}
			if len(parts) > 7 {
				symbolGlyphs[parts[2]] = parts[6]
				symbolGlyphs[parts[5]] = parts[7]
			}
			primarySymbolMap[l[0]] = parts[2]
			alternateSymbolMap[l[0]] = parts[5]
		}
	}
}
