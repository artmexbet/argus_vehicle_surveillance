--
-- PostgreSQL database dump
--

-- Dumped from database version 15.8 (Debian 15.8-1.pgdg120+1)
-- Dumped by pg_dump version 16.4

-- Started on 2024-11-02 21:45:47

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 5 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- TOC entry 3418 (class 0 OID 0)
-- Dependencies: 5
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 233 (class 1255 OID 24601)
-- Name: on_account_created(); Type: FUNCTION; Schema: public; Owner: baseuser
--

CREATE FUNCTION public.on_account_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	INSERT INTO public.profile 
	VALUES (NEW.id, '', '', '');
	RETURN NEW;
END;$$;


ALTER FUNCTION public.on_account_created() OWNER TO baseuser;

--
-- TOC entry 234 (class 1255 OID 24621)
-- Name: on_car_signalized(); Type: FUNCTION; Schema: public; Owner: baseuser
--

CREATE FUNCTION public.on_car_signalized() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	UPDATE public.accounts
	SET cars_under_shield = (
		SELECT coalesce(json_agg(row_to_json(t)), '[]') FROM
			(SELECT * FROM public.security_cars
			WHERE account_id=NEW.account_id AND security_date_off IS NULL
			ORDER BY security_date_on) AS t)
	WHERE id = NEW.account_id;
	RETURN NEW;
END;$$;


ALTER FUNCTION public.on_car_signalized() OWNER TO baseuser;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 16505)
-- Name: accounts; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.accounts (
    id uuid DEFAULT public.uuid_generate_v1() NOT NULL,
    login character varying(50) NOT NULL,
    password character varying(200) NOT NULL,
    cars_under_shield json
);


ALTER TABLE public.accounts OWNER TO baseuser;

--
-- TOC entry 218 (class 1259 OID 16534)
-- Name: cameras; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.cameras (
    id uuid DEFAULT public.uuid_generate_v1() NOT NULL,
    ip_address character varying(15) NOT NULL
);


ALTER TABLE public.cameras OWNER TO baseuser;

--
-- TOC entry 221 (class 1259 OID 24578)
-- Name: event_types; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.event_types (
    id integer NOT NULL,
    title character varying(50) NOT NULL
);


ALTER TABLE public.event_types OWNER TO baseuser;

--
-- TOC entry 220 (class 1259 OID 24577)
-- Name: event_types_id_seq; Type: SEQUENCE; Schema: public; Owner: baseuser
--

CREATE SEQUENCE public.event_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.event_types_id_seq OWNER TO baseuser;

--
-- TOC entry 3419 (class 0 OID 0)
-- Dependencies: 220
-- Name: event_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: baseuser
--

ALTER SEQUENCE public.event_types_id_seq OWNED BY public.event_types.id;


--
-- TOC entry 222 (class 1259 OID 24586)
-- Name: events; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.events (
    id uuid DEFAULT public.uuid_generate_v1() NOT NULL,
    sc_id uuid NOT NULL,
    et_id bigint NOT NULL,
    "time" timestamp with time zone NOT NULL
);


ALTER TABLE public.events OWNER TO baseuser;

--
-- TOC entry 216 (class 1259 OID 16512)
-- Name: profile; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.profile (
    account_id uuid DEFAULT public.uuid_generate_v1() NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    phone_number character varying(20)
);


ALTER TABLE public.profile OWNER TO baseuser;

--
-- TOC entry 219 (class 1259 OID 16539)
-- Name: security_cars; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.security_cars (
    id uuid DEFAULT public.uuid_generate_v1() NOT NULL,
    camera_id uuid NOT NULL,
    account_id uuid NOT NULL,
    security_date_on timestamp with time zone NOT NULL,
    security_date_off timestamp with time zone,
    car_id bigint
);


ALTER TABLE public.security_cars OWNER TO baseuser;

--
-- TOC entry 217 (class 1259 OID 16522)
-- Name: tokens; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.tokens (
    id uuid DEFAULT public.uuid_generate_v1() NOT NULL,
    account_id uuid NOT NULL,
    create_date timestamp with time zone NOT NULL,
    refresh bytea NOT NULL
);


ALTER TABLE public.tokens OWNER TO baseuser;

--
-- TOC entry 3241 (class 2604 OID 24581)
-- Name: event_types id; Type: DEFAULT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.event_types ALTER COLUMN id SET DEFAULT nextval('public.event_types_id_seq'::regclass);


--
-- TOC entry 3244 (class 2606 OID 16557)
-- Name: accounts accounts_login_key; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_login_key UNIQUE (login);


--
-- TOC entry 3246 (class 2606 OID 16511)
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- TOC entry 3252 (class 2606 OID 16555)
-- Name: cameras cameras_ip_address_key; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.cameras
    ADD CONSTRAINT cameras_ip_address_key UNIQUE (ip_address);


--
-- TOC entry 3254 (class 2606 OID 16538)
-- Name: cameras cameras_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.cameras
    ADD CONSTRAINT cameras_pkey PRIMARY KEY (id);


--
-- TOC entry 3258 (class 2606 OID 24583)
-- Name: event_types event_types_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.event_types
    ADD CONSTRAINT event_types_pkey PRIMARY KEY (id);


--
-- TOC entry 3260 (class 2606 OID 24585)
-- Name: event_types event_types_title_key; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.event_types
    ADD CONSTRAINT event_types_title_key UNIQUE (title);


--
-- TOC entry 3262 (class 2606 OID 24590)
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- TOC entry 3248 (class 2606 OID 16516)
-- Name: profile profile_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.profile
    ADD CONSTRAINT profile_pkey PRIMARY KEY (account_id);


--
-- TOC entry 3256 (class 2606 OID 16543)
-- Name: security_cars security_cars_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.security_cars
    ADD CONSTRAINT security_cars_pkey PRIMARY KEY (id);


--
-- TOC entry 3250 (class 2606 OID 16528)
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 3269 (class 2620 OID 24603)
-- Name: accounts on_account_created; Type: TRIGGER; Schema: public; Owner: baseuser
--

CREATE TRIGGER on_account_created AFTER INSERT ON public.accounts FOR EACH ROW EXECUTE FUNCTION public.on_account_created();


--
-- TOC entry 3270 (class 2620 OID 24625)
-- Name: security_cars on_insert_update_delete; Type: TRIGGER; Schema: public; Owner: baseuser
--

CREATE TRIGGER on_insert_update_delete AFTER INSERT OR DELETE OR UPDATE ON public.security_cars FOR EACH ROW EXECUTE FUNCTION public.on_car_signalized();


--
-- TOC entry 3267 (class 2606 OID 24596)
-- Name: events events_et_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_et_id_fkey FOREIGN KEY (et_id) REFERENCES public.event_types(id);


--
-- TOC entry 3268 (class 2606 OID 24591)
-- Name: events events_sc_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_sc_id_fkey FOREIGN KEY (sc_id) REFERENCES public.security_cars(id);


--
-- TOC entry 3263 (class 2606 OID 16517)
-- Name: profile profile_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.profile
    ADD CONSTRAINT profile_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id);


--
-- TOC entry 3265 (class 2606 OID 16549)
-- Name: security_cars security_cars_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.security_cars
    ADD CONSTRAINT security_cars_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id);


--
-- TOC entry 3266 (class 2606 OID 16544)
-- Name: security_cars security_cars_camera_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.security_cars
    ADD CONSTRAINT security_cars_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.cameras(id);


--
-- TOC entry 3264 (class 2606 OID 16529)
-- Name: tokens tokens_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id);


-- Completed on 2024-11-02 21:45:48

--
-- PostgreSQL database dump complete
--

